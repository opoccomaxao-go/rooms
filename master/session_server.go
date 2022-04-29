package master

import (
	"log"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/ipc/event"
	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/pkg/errors"
)

type SessionServer struct {
	id      uint64
	conn    *channel.Channel
	handler ClientHandler
	logger  *log.Logger
}

type ClientHandler interface {
	Register(id uint64, client *SessionServer)
}

func (s *SessionServer) Handle(event *event.Common) {

	panic("implement")
}

func (s *SessionServer) Serve() {
	err := errors.WithStack(s.conn.Serve(s))
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
