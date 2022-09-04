package session

import (
	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/ipc/event"
	"github.com/opoccomaxao-go/ipc/processor"
	"github.com/opoccomaxao-go/rooms/apm"
	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/opoccomaxao-go/rooms/proto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type connWrapper struct {
	conn     *channel.Client
	parent   *Server
	logger   zerolog.Logger
	interval apm.DebuggableInterval
}

func (c *connWrapper) init() {
	c.logger = c.parent.config.Logger.With().Logger()
	c.interval = apm.NewZerologInterval(&c.logger, "session.connWrapper.")
}

func (c *connWrapper) Handler() channel.Handler[*event.Common] {
	defer c.interval.Start("Handler").End()

	res := processor.New()

	res.Register(proto.CommandMasterAuthRequired, c.onAuthRequired)
	res.Register(proto.CommandMasterAuthSuccess, c.onAuthSuccess)
	res.Register(proto.CommandMasterRoomCreate, c.onRoomCreate)
	res.Register(proto.CommandMasterRoomCancel, c.onRoomCancel)

	return res
}

func (c *connWrapper) Serve() error {
	defer c.interval.Start("Serve").End()

	return c.conn.Serve(c.Handler())
}

func (c *connWrapper) onAuthRequired(payload []byte) {
	defer c.interval.Start("onAuthRequired").End()

	if len(payload) == 0 {
		c.Auth(&proto.Auth{
			Version: constants.Version,
			Token:   string(c.parent.config.Token),
		})
	} else {
		c.parent.onAuthError(string(payload))
	}
}

func (c *connWrapper) onAuthSuccess(_ []byte) {
	defer c.interval.Start("onAuthSuccess").End()

	c.Stats(&proto.Stats{
		Capacity: c.parent.getCapacity(),
	})
}

func (c *connWrapper) onRoomCreate(payload []byte) {
	defer c.interval.Start("onRoomCreate").End()

	var room proto.Room

	err := errors.WithStack(room.Read(payload))
	if err != nil {
		c.logger.Err(err).Stack().Send()
	}

	c.parent.onRoomCreate(&room)
}

func (c *connWrapper) onRoomCancel(payload []byte) {
	defer c.interval.Start("onRoomCancel").End()

	c.parent.onRoomCancel(proto.ReadID(payload))
}

func (c *connWrapper) Auth(auth *proto.Auth) {
	defer c.interval.Start("Auth").End()

	c.conn.Send(&event.Common{
		Type:    proto.CommandSessionAuth,
		Payload: auth.Payload(),
	})
}

func (c *connWrapper) RoomCreated(room *proto.Room) {
	defer c.interval.Start("RoomCreated").End()

	c.conn.Send(&event.Common{
		Type:    proto.CommandSessionRoomCreated,
		Payload: room.Payload(),
	})
}

func (c *connWrapper) RoomError(room *proto.Room) {
	defer c.interval.Start("RoomError").End()

	c.conn.Send(&event.Common{
		Type:    proto.CommandSessionRoomError,
		Payload: room.Payload(),
	})
}

func (c *connWrapper) RoomFinished(room *proto.Room) {
	defer c.interval.Start("RoomFinished").End()

	c.conn.Send(&event.Common{
		Type:    proto.CommandSessionRoomFinished,
		Payload: room.Payload(),
	})
}

func (c *connWrapper) Stats(stats *proto.Stats) {
	defer c.interval.Start("Stats").End()

	c.conn.Send(&event.Common{
		Type:    proto.CommandSessionStats,
		Payload: stats.Payload(),
	})
}

func (c *connWrapper) Close() error {
	defer c.interval.Start("Close").End()

	err := c.conn.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
