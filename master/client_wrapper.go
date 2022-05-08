package master

import (
	"context"
	"log"
	"sync"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/ipc/event"
	"github.com/opoccomaxao-go/ipc/processor"
	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/opoccomaxao-go/rooms/proto"
	"github.com/opoccomaxao-go/rooms/utils"
	"github.com/pkg/errors"
)

type RoomCreateResult struct {
	Room  *proto.Room
	Error error
}

type clientWrapper struct {
	conn   *channel.Channel
	parent *Server
	logger *log.Logger

	// internal set only

	id        uint64
	stats     proto.Stats
	listeners map[proto.ID][]chan RoomCreateResult

	mu sync.RWMutex
}

func (s *clientWrapper) onAuth(payload []byte) {
	id, err := s.parent.config.Storage.Validate(payload)
	if err != nil {
		s.logger.Printf("%v\n", errors.WithStack(err))
		s.AuthRequired(err)

		return
	}

	s.parent.register(id, s)
	s.parent.unregister(s.id, s)
	s.id = id
	s.AuthSuccess()
}

func (s *clientWrapper) onRoomCreated(payload []byte) {
	var room proto.Room

	err := room.Read(payload)
	if err != nil {
		s.logger.Printf("%v\n", errors.WithStack(err))

		return
	}

	s.notifyRoomCreate(room.ID, RoomCreateResult{
		Room: &room,
	})
}

func (s *clientWrapper) onRoomError(payload []byte) {
	var room proto.Room

	err := room.Read(payload)
	if err != nil {
		s.logger.Printf("%v\n", errors.WithStack(err))

		return
	}

	s.notifyRoomCreate(room.ID, RoomCreateResult{
		Error: errors.New(room.Error),
	})
}

func (s *clientWrapper) onRoomFinished(payload []byte) {
	var room proto.Room

	err := room.Read(payload)
	if err != nil {
		s.logger.Printf("%v\n", errors.WithStack(err))

		return
	}

	s.parent.notifyFinishedRoom(&room)
}

func (s *clientWrapper) onStats(payload []byte) {
	err := errors.WithStack(s.stats.Read(payload))
	if err != nil {
		s.logger.Printf("%v\n", err)
	}

	s.parent.onStats()
}

func (s *clientWrapper) Serve() {
	s.clearWaiters()

	handler := processor.New()
	handler.Register(proto.CommandSessionAuth, s.onAuth)
	handler.Register(proto.CommandSessionRoomCreated, s.onRoomCreated)
	handler.Register(proto.CommandSessionRoomError, s.onRoomError)
	handler.Register(proto.CommandSessionRoomFinished, s.onRoomFinished)
	handler.Register(proto.CommandSessionStats, s.onStats)

	s.AuthRequired(nil)

	err := errors.WithStack(s.conn.Serve(handler))
	if err != nil {
		s.logger.Printf("%v\n", err)
	}
}

// FlushInstance take all unsent data from other equal server.
func (s *clientWrapper) FlushInstance(other *clientWrapper) error {
	if s.id != other.id {
		return errors.Wrapf(constants.ErrInvalid, "illegal instance id: %d, required %d", other.id, s.id)
	}

	err := other.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *clientWrapper) addWaiter(id uint64, waiter chan RoomCreateResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.listeners[id] = append(s.listeners[id], waiter)
}

func (s *clientWrapper) removeWaiters(id uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.listeners, id)
}

func (s *clientWrapper) notifyRoomCreate(id uint64, result RoomCreateResult) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	utils.WithChannels(s.listeners[id]).Notify(result)
}

func (s *clientWrapper) clearWaiters() {
	s.mu.Lock()
	defer s.mu.Unlock()

	go s.closeAll(s.listeners)

	s.listeners = map[uint64][]chan RoomCreateResult{}
}

func (*clientWrapper) closeAll(allWaiters map[uint64][]chan RoomCreateResult) {
	if allWaiters == nil {
		return
	}

	for _, waiters := range allWaiters {
		utils.WithChannels(waiters).Close()
	}
}

func (s *clientWrapper) WaitRoomCreateResult(ctx context.Context, id uint64) <-chan RoomCreateResult {
	waiter := make(chan RoomCreateResult, 1)

	s.addWaiter(id, waiter)

	utils.WithChannel(waiter).
		BeforeClose(func() { s.removeWaiters(id) }).
		AsyncCloseOnDone(ctx)

	return waiter
}

func (s *clientWrapper) AuthRequired(err error) {
	event := event.Common{
		Type: proto.CommandMasterAuthRequired,
	}

	if err != nil {
		event.Payload = []byte(err.Error())
	}

	s.conn.Send(&event)
}

func (s *clientWrapper) AuthSuccess() {
	s.conn.Send(&event.Common{
		Type: proto.CommandMasterAuthSuccess,
	})
}

func (s *clientWrapper) RoomCreate(room *proto.Room) {
	s.conn.Send(&event.Common{
		Type:    proto.CommandMasterRoomCreate,
		Payload: room.Payload(),
	})
}

func (s *clientWrapper) RoomCancel(roomID proto.ID) {
	s.conn.Send(&event.Common{
		Type:    proto.CommandMasterRoomCancel,
		Payload: proto.PayloadID(roomID),
	})
}

func (s *clientWrapper) Close() error {
	err := s.conn.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	s.clearWaiters()

	return nil
}
