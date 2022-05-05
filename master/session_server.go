package master

import (
	"log"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/ipc/event"
	"github.com/opoccomaxao-go/ipc/processor"
	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/opoccomaxao-go/rooms/proto"
	"github.com/pkg/errors"
)

type SessionServer struct {
	id     uint64
	conn   *channel.Channel
	parent *SessionStorage
	logger *log.Logger
	Stats  proto.Stats
}

func (s *SessionServer) auth(payload []byte) {
	id, err := s.parent.storage.Validate(payload)
	if err != nil {
		s.logger.Printf("%v\n", errors.WithStack(err))
		s.conn.Send(&event.Common{
			Type:    proto.CommandMasterAuthRequired,
			Payload: []byte(err.Error()),
		})
		return
	}

	s.id = id
	s.conn.Send(&event.Common{
		Type: proto.CommandMasterAuthSuccess,
	})
}

func (s *SessionServer) roomCreated(payload []byte) {
	panic("implement")
}

func (s *SessionServer) roomError(payload []byte) {
	panic("implement")
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
	handler := processor.New()
	handler.Register(proto.CommandSessionAuth, s.auth)
	handler.Register(proto.CommandSessionRoomCreated, s.roomCreated)
	handler.Register(proto.CommandSessionRoomError, s.roomError)
	handler.Register(proto.CommandSessionRoomFinished, s.roomFinished)
	handler.Register(proto.CommandSessionStats, s.stats)

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
