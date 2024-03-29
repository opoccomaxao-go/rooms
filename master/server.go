package master

import (
	"context"
	"sync"
	"time"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/rooms/apm"
	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/opoccomaxao-go/rooms/proto"
	"github.com/opoccomaxao-go/rooms/storage"
	"github.com/opoccomaxao-go/rooms/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/exp/slices"
)

const DefaultRoomListenerCapacity = 100

type Server struct {
	config   Config
	interval apm.DebuggableInterval

	server  *channel.Server
	clients map[uint64]*connWrapper

	condStats         *sync.Cond
	listenersFinished []chan *proto.Room

	mu sync.RWMutex
}

type Config struct {
	Logger         *zerolog.Logger
	Storage        storage.Storage
	SessionAddress string        // SessionAddress is address for session-server listening.
	CreateTimeout  time.Duration // CreateTimeout is NewRoom timeout.
}

func New(cfg Config) (*Server, error) {
	var err error

	if cfg.Logger == nil {
		logger := zerolog.Nop()
		cfg.Logger = &logger
	}

	if cfg.Storage == nil {
		cfg.Storage = storage.NewRAM()
	}

	if cfg.SessionAddress == "" {
		cfg.SessionAddress = constants.DefaultAddress
	}

	if cfg.CreateTimeout <= 0 {
		cfg.CreateTimeout = constants.DefaultTimeout
	}

	res := &Server{
		config:    cfg,
		interval:  apm.NewZerologInterval(cfg.Logger, "master.Server."),
		clients:   map[uint64]*connWrapper{},
		condStats: sync.NewCond(&sync.Mutex{}),
	}

	res.server, err = channel.NewServer(channel.ServerConfig{
		Address: cfg.SessionAddress,
		Handler: channel.HandlerFunc[*channel.Channel](res.handle),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}

func (s *Server) handle(conn *channel.Channel) {
	defer s.interval.Start("handle").End()

	server := connWrapper{
		conn:   conn,
		parent: s,
	}

	server.init()

	server.Serve()
}

func (s *Server) register(id uint64, client *connWrapper) {
	defer s.interval.Start("register").End()

	s.mu.Lock()
	defer s.mu.Unlock()

	if prev, ok := s.clients[id]; ok {
		err := errors.Wrap(client.FlushInstance(prev), "flush error")
		if err != nil {
			s.config.Logger.Err(err).Stack().Send()
		}
	}

	s.clients[id] = client
}

func (s *Server) unregister(id uint64, client *connWrapper) {
	defer s.interval.Start("unregister").End()

	s.mu.Lock()
	defer s.mu.Unlock()

	if prev, ok := s.clients[id]; ok {
		if prev == client {
			delete(s.clients, id)
		}
	}
}

func (s *Server) Serve(ctx context.Context) error {
	defer s.interval.Start("Serve").End()

	utils.WithContext(ctx).
		AsyncOnDone(func() {
			err := s.Close()
			if err != nil {
				s.config.Logger.Err(err).Stack().Send()
			}
		})

	return errors.WithStack(s.server.Listen())
}

func (s *Server) Close() error {
	defer s.interval.Start("Close").End()

	return errors.WithStack(s.server.Close())
}

func (s *Server) findFreeServer() *connWrapper {
	defer s.interval.Start("findFreeServer").End()

	var best *connWrapper

	for _, ss := range s.clients {
		if ss.stats.Capacity > 0 && (best == nil || ss.stats.Capacity > best.stats.Capacity) {
			best = ss
		}
	}

	return best
}

func (s *Server) onStats() {
	defer s.interval.Start("onStats").End()

	s.condStats.Broadcast()
}

func (s *Server) waitStats() <-chan struct{} {
	defer s.interval.Start("waitStats").End()

	res := make(chan struct{}, 1)

	utils.WithChannel(res).
		BeforeClose(func() { res <- struct{}{} }).
		AsyncCloseAfterFunc(s.condStats.Wait)

	return res
}

func (s *Server) CreateRoom(ctx context.Context, userIDs []uint64) (*proto.Room, error) {
	defer s.interval.Start("CreateRoom").End()

	s.mu.Lock()
	defer s.mu.Unlock()

	room := &proto.Room{
		ID:      s.config.Storage.NewRoom(),
		Clients: make([]*proto.Client, len(userIDs)),
	}

	for index, id := range userIDs {
		room.Clients[index] = &proto.Client{
			ID: id,
		}
	}

	ctx, cancelFn := context.WithTimeout(ctx, s.config.CreateTimeout)
	defer cancelFn()

	done := ctx.Done()

	for {
		select {
		case <-done:
			return nil, ctx.Err()
		default:
		}

		best := s.findFreeServer()

		if best == nil {
			select {
			case <-done:
				return nil, ctx.Err()
			case <-s.waitStats():
				continue
			}
		}

		waiter := best.WaitRoomCreateResult(ctx, room.ID)

		best.RoomCreate(room)

		res := <-waiter

		if res.Error != nil {
			best.RoomCancel(room.ID)
		}

		if res.Room == nil {
			continue
		}

		res.Room.ServerID = best.id

		return res.Room, res.Error
	}
}

func (s *Server) pushFinishedListener(listener chan *proto.Room) {
	defer s.interval.Start("pushFinishedListener").End()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.listenersFinished = append(s.listenersFinished, listener)
}

func (s *Server) removeFinishedListener(listener chan *proto.Room) {
	defer s.interval.Start("removeFinishedListener").End()

	s.mu.Lock()
	defer s.mu.Unlock()

	index := slices.Index(s.listenersFinished, listener)
	if index == -1 {
		return
	}

	s.listenersFinished = slices.Delete(s.listenersFinished, index, 1)
}

// FinishedRooms creates channel-receiver of all finished rooms. To close channel cancel context ctx.
func (s *Server) FinishedRooms(ctx context.Context) <-chan *proto.Room {
	defer s.interval.Start("FinishedRooms").End()

	res := make(chan *proto.Room, DefaultRoomListenerCapacity)

	s.pushFinishedListener(res)

	utils.WithChannel(res).
		BeforeClose(func() { s.removeFinishedListener(res) }).
		AsyncCloseOnDone(ctx)

	return res
}

func (s *Server) notifyFinishedRoom(room *proto.Room) {
	defer s.interval.Start("notifyFinishedRoom").End()

	s.mu.RLock()
	defer s.mu.RUnlock()

	utils.WithChannels(s.listenersFinished).Notify(room)
}
