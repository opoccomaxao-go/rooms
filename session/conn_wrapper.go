package session

import (
	"log"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/ipc/event"
	"github.com/opoccomaxao-go/ipc/processor"
	"github.com/opoccomaxao-go/rooms/proto"
	"github.com/pkg/errors"
)

type connWrapper struct {
	conn   *channel.Client
	parent *Server
	logger *log.Logger
}

func (c *connWrapper) Handler() channel.Handler[*event.Common] {
	res := processor.New()

	res.Register(proto.CommandMasterAuthRequired, c.onAuthRequired)
	res.Register(proto.CommandMasterAuthSuccess, c.onAuthSuccess)
	res.Register(proto.CommandMasterRoomCreate, c.onRoomCreate)
	res.Register(proto.CommandMasterRoomCancel, c.onRoomCancel)

	return res
}

func (c *connWrapper) onAuthRequired(payload []byte) {
	// TODO: implement
}

func (c *connWrapper) onAuthSuccess(payload []byte) {
	// TODO: implement
}

func (c *connWrapper) onRoomCreate(payload []byte) {
	// TODO: implement
}

func (c *connWrapper) onRoomCancel(payload []byte) {
	// TODO: implement
}

func (c *connWrapper) Auth(auth *proto.Auth) {
	c.conn.Send(&event.Common{
		Type:    proto.CommandSessionAuth,
		Payload: auth.Payload(),
	})
}

func (c *connWrapper) RoomCreated(room *proto.Room) {
	c.conn.Send(&event.Common{
		Type:    proto.CommandSessionRoomCreated,
		Payload: room.Payload(),
	})
}

func (c *connWrapper) RoomError(room *proto.Room) {
	c.conn.Send(&event.Common{
		Type:    proto.CommandSessionRoomError,
		Payload: room.Payload(),
	})
}

func (c *connWrapper) RoomFinished(room *proto.Room) {
	c.conn.Send(&event.Common{
		Type:    proto.CommandSessionRoomFinished,
		Payload: room.Payload(),
	})
}

func (c *connWrapper) Stats(stats *proto.Stats) {
	c.conn.Send(&event.Common{
		Type:    proto.CommandSessionStats,
		Payload: stats.Payload(),
	})
}

func (c *connWrapper) Close() error {
	err := c.conn.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
