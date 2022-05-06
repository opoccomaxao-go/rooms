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
	"github.com/pkg/errors"
)

type RoomCreateResult struct {
	Room  *proto.Room
	Error error
}

type SessionServer struct {
	id     uint64
	conn   *channel.Channel
	parent *SessionStorage
	logger *log.Logger
	Stats  proto.Stats

	waiters map[proto.ID][]chan RoomCreateResult
	mu      sync.RWMutex
}

func (s *SessionServer) auth(payload []byte) {
	id, err := s.parent.storage.Validate(payload)
	if err != nil {
		s.logger.Printf("%v\n", errors.WithStack(err))
		s.AuthRequired(err)

		return
	}

	s.id = id
	s.AuthSuccess()
}

func (s *SessionServer) roomCreated(payload []byte) {
	var room proto.Room

	err := room.Read(payload)
	if err != nil {
		s.logger.Printf("%v\n", errors.WithStack(err))

		return
	}

	s.notifyRoomCreate(room.ID, RoomCreateResult{
		Room:  &room,
		Error: nil,
	})
}

func (s *SessionServer) roomError(payload []byte) {
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

func (s *SessionServer) roomFinished(payload []byte) {
	panic("implement")
}

func (s *SessionServer) stats(payload []byte) {
	err := errors.WithStack(s.Stats.Read(payload))
	if err != nil {
		s.logger.Printf("%v\n", err)
	}

	s.parent.TriggerStats()
}

func (s *SessionServer) Serve() {
	s.clearWaiters()

	handler := processor.New()
	handler.Register(proto.CommandSessionAuth, s.auth)
	handler.Register(proto.CommandSessionRoomCreated, s.roomCreated)
	handler.Register(proto.CommandSessionRoomError, s.roomError)
	handler.Register(proto.CommandSessionRoomFinished, s.roomFinished)
	handler.Register(proto.CommandSessionStats, s.stats)

	s.AuthRequired(nil)

	err := errors.WithStack(s.conn.Serve(handler))
	if err != nil {
		s.logger.Printf("%v\n", err)
	}
}

// FlushInstance take all unsent data from other equal server.
func (s *SessionServer) FlushInstance(other *SessionServer) error {
	if s.id != other.id {
		return errors.Wrapf(constants.ErrInvalid, "illegal instance id: %d, required %d", other.id, s.id)
	}

	panic("implement")
}

func (s *SessionServer) addWaiter(id uint64, waiter chan RoomCreateResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.waiters[id] = append(s.waiters[id], waiter)
}

func (s *SessionServer) removeWaiters(id uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.waiters, id)
}

func (s *SessionServer) notifyRoomCreate(id uint64, result RoomCreateResult) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, waiter := range s.waiters[id] {
		waiter <- result
	}
}

func (s *SessionServer) clearWaiters() {
	s.mu.Lock()
	defer s.mu.Unlock()

	go s.closeAll(s.waiters)

	s.waiters = map[uint64][]chan RoomCreateResult{}
}

func (*SessionServer) closeAll(allWaiters map[uint64][]chan RoomCreateResult) {
	if allWaiters == nil {
		return
	}

	for _, waiters := range allWaiters {
		for _, waiter := range waiters {
			close(waiter)
		}
	}
}

func (s *SessionServer) WaitRoomCreateResult(ctx context.Context, id uint64) <-chan RoomCreateResult {
	waiter := make(chan RoomCreateResult, 1)

	s.addWaiter(id, waiter)

	go func() {
		<-ctx.Done()
		s.removeWaiters(id)
		close(waiter)
	}()

	return waiter
}

func (s *SessionServer) AuthRequired(err error) {
	event := event.Common{
		Type: proto.CommandMasterAuthRequired,
	}

	if err != nil {
		event.Payload = []byte(err.Error())
	}

	s.conn.Send(&event)
}

func (s *SessionServer) AuthSuccess() {
	s.conn.Send(&event.Common{
		Type: proto.CommandMasterAuthSuccess,
	})
}

func (s *SessionServer) RoomCreate(room *proto.Room) {
	s.conn.Send(&event.Common{
		Type:    proto.CommandMasterRoomCreate,
		Payload: room.Payload(),
	})
}

func (s *SessionServer) RoomCancel(roomID proto.ID) {
	s.conn.Send(&event.Common{
		Type:    proto.CommandMasterRoomCancel,
		Payload: proto.PayloadID(roomID),
	})
}
