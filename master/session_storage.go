package master

import (
	"context"
	"log"
	"sync"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/rooms/proto"
	"github.com/opoccomaxao-go/rooms/storage"
	"github.com/pkg/errors"
)

type SessionStorage struct {
	clients map[uint64]*SessionServer
	logger  *log.Logger
	storage storage.Storage
	trigger *sync.Cond
	mu      sync.Mutex
}

func newSessionStorage(cfg *Config) *SessionStorage {
	return &SessionStorage{
		clients: map[uint64]*SessionServer{},
		logger:  cfg.Logger,
		storage: cfg.Storage,
		trigger: sync.NewCond(&sync.Mutex{}),
	}
}

func (s *SessionStorage) Handle(conn *channel.Channel) {
	server := SessionServer{
		conn:   conn,
		parent: s,
		logger: s.logger,
	}

	server.Serve()
}

func (s *SessionStorage) Register(id uint64, client *SessionServer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if prev, ok := s.clients[id]; ok {
		err := errors.Wrap(client.FlushInstance(prev), "flush error")
		if err != nil {
			s.logger.Printf("%v\n", err)
		}
	}

	s.clients[id] = client
}

func (s *SessionStorage) Unregister(id uint64, client *SessionServer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if prev, ok := s.clients[id]; ok {
		if prev == client {
			delete(s.clients, id)
		}
	}
}

func (s *SessionStorage) TriggerStats() {
	s.trigger.Broadcast()
}

func (s *SessionStorage) WaitStats() <-chan struct{} {
	res := make(chan struct{}, 1)

	go func() {
		s.trigger.Wait()
		res <- struct{}{}
		close(res)
	}()

	return res
}

func (s *SessionStorage) getFreeServer() *SessionServer {
	s.mu.Lock()
	defer s.mu.Unlock()

	var best *SessionServer

	for _, ss := range s.clients {
		if ss.Stats.Capacity > 0 && (best == nil || ss.Stats.Capacity > best.Stats.Capacity) {
			best = ss
		}
	}

	return best
}

func (s *SessionStorage) CreateRoom(ctx context.Context, room *proto.Room) (*proto.Room, error) {
	done := ctx.Done()

	for {
		best := s.getFreeServer()

		if best == nil {
			select {
			case <-done:
				return nil, ctx.Err()
			case <-s.WaitStats():
				continue
			}
		}

		waiter := best.WaitRoomCreateResult(ctx, room.ID)

		best.RoomCreate(room)

		res := <-waiter

		if res.Error != nil {
			best.RoomCancel(room.ID)
		}

		return res.Room, res.Error
	}
}
