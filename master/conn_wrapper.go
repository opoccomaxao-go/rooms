package master

import (
	"context"
	"sync"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/ipc/event"
	"github.com/opoccomaxao-go/ipc/processor"
	"github.com/opoccomaxao-go/rooms/apm"
	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/opoccomaxao-go/rooms/proto"
	"github.com/opoccomaxao-go/rooms/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type RoomCreateResult struct {
	Room  *proto.Room
	Error error
}

type connWrapper struct {
	conn     *channel.Channel
	parent   *Server
	logger   zerolog.Logger
	interval apm.DebuggableInterval

	// internal set only

	id        uint64
	stats     proto.Stats
	listeners map[proto.ID][]chan RoomCreateResult

	mu sync.RWMutex
}

func (c *connWrapper) init() {
	c.logger = c.parent.config.Logger.With().
		Logger()
	c.interval = apm.NewZerologInterval(&c.logger, "master.connWrapper.")
	c.listeners = map[uint64][]chan RoomCreateResult{}
}

func (c *connWrapper) onAuth(payload []byte) {
	defer c.interval.Start("onAuth").End()

	var auth proto.Auth

	err := auth.Read(payload)
	if err != nil {
		c.logger.Err(err).Stack().Send()

		c.AuthRequired(err)

		return
	}

	id, err := c.parent.config.Storage.Validate(auth.Version, auth.Token)
	if err != nil {
		c.logger.Err(err).Stack().Send()

		c.AuthRequired(err)

		return
	}

	c.parent.register(id, c)
	c.parent.unregister(c.id, c)
	c.id = id
	c.AuthSuccess()
}

func (c *connWrapper) onRoomCreated(payload []byte) {
	defer c.interval.Start("onRoomCreated").End()

	var room proto.Room

	err := room.Read(payload)
	if err != nil {
		c.logger.Err(err).Stack().Send()

		return
	}

	c.notifyRoomCreate(room.ID, RoomCreateResult{
		Room: &room,
	})
}

func (c *connWrapper) onRoomError(payload []byte) {
	defer c.interval.Start("onRoomError").End()

	var room proto.Room

	err := room.Read(payload)
	if err != nil {
		c.logger.Err(err).Stack().Send()

		return
	}

	c.notifyRoomCreate(room.ID, RoomCreateResult{
		Error: errors.New(room.Error),
	})
}

func (c *connWrapper) onRoomFinished(payload []byte) {
	defer c.interval.Start("onRoomFinished").End()

	var room proto.Room

	err := room.Read(payload)
	if err != nil {
		c.logger.Err(err).Stack().Send()

		return
	}

	c.parent.notifyFinishedRoom(&room)
}

func (c *connWrapper) onStats(payload []byte) {
	defer c.interval.Start("onStats").End()

	err := errors.WithStack(c.stats.Read(payload))
	if err != nil {
		c.logger.Err(err).Stack().Send()
	}

	c.parent.onStats()
}

func (c *connWrapper) Serve() {
	defer c.interval.Start("Serve").End()

	c.clearWaiters()

	handler := processor.New()
	handler.Register(proto.CommandSessionAuth, c.onAuth)
	handler.Register(proto.CommandSessionRoomCreated, c.onRoomCreated)
	handler.Register(proto.CommandSessionRoomError, c.onRoomError)
	handler.Register(proto.CommandSessionRoomFinished, c.onRoomFinished)
	handler.Register(proto.CommandSessionStats, c.onStats)

	c.AuthRequired(nil)

	err := errors.WithStack(c.conn.Serve(handler))
	if err != nil {
		c.logger.Err(err).Stack().Send()
	}
}

// FlushInstance take all unsent data from other equal server.
func (c *connWrapper) FlushInstance(other *connWrapper) error {
	defer c.interval.Start("FlushInstance").End()

	if c.id != other.id {
		return errors.Wrapf(constants.ErrInvalid, "illegal instance id: %d, required %d", other.id, c.id)
	}

	err := other.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *connWrapper) addWaiter(id uint64, waiter chan RoomCreateResult) {
	defer c.interval.Start("addWaiter").End()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.listeners[id] = append(c.listeners[id], waiter)
}

func (c *connWrapper) removeWaiters(id uint64) {
	defer c.interval.Start("removeWaiters").End()

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.listeners, id)
}

func (c *connWrapper) notifyRoomCreate(id uint64, result RoomCreateResult) {
	defer c.interval.Start("notifyRoomCreate").End()

	c.mu.RLock()
	defer c.mu.RUnlock()

	utils.WithChannels(c.listeners[id]).Notify(result)
}

func (c *connWrapper) clearWaiters() {
	defer c.interval.Start("clearWaiters").End()

	c.mu.Lock()
	defer c.mu.Unlock()

	go c.closeAll(c.listeners)

	c.listeners = map[uint64][]chan RoomCreateResult{}
}

func (c *connWrapper) closeAll(allWaiters map[uint64][]chan RoomCreateResult) {
	defer c.interval.Start("closeAll").End()

	if allWaiters == nil {
		return
	}

	for _, waiters := range allWaiters {
		utils.WithChannels(waiters).Close()
	}
}

func (c *connWrapper) WaitRoomCreateResult(ctx context.Context, id uint64) <-chan RoomCreateResult {
	defer c.interval.Start("WaitRoomCreateResult").End()

	waiter := make(chan RoomCreateResult, 1)

	c.addWaiter(id, waiter)

	utils.WithChannel(waiter).
		BeforeClose(func() { c.removeWaiters(id) }).
		AsyncCloseOnDone(ctx)

	return waiter
}

func (c *connWrapper) AuthRequired(err error) {
	defer c.interval.Start("AuthRequired").End()

	event := event.Common{
		Type: proto.CommandMasterAuthRequired,
	}

	if err != nil {
		event.Payload = []byte(err.Error())
	}

	c.conn.Send(&event)
}

func (c *connWrapper) AuthSuccess() {
	defer c.interval.Start("AuthSuccess").End()

	c.conn.Send(&event.Common{
		Type: proto.CommandMasterAuthSuccess,
	})
}

func (c *connWrapper) RoomCreate(room *proto.Room) {
	defer c.interval.Start("RoomCreate").End()

	c.conn.Send(&event.Common{
		Type:    proto.CommandMasterRoomCreate,
		Payload: room.Payload(),
	})
}

func (c *connWrapper) RoomCancel(roomID proto.ID) {
	defer c.interval.Start("RoomCancel").End()

	c.conn.Send(&event.Common{
		Type:    proto.CommandMasterRoomCancel,
		Payload: proto.PayloadID(roomID),
	})
}

func (c *connWrapper) Close() error {
	defer c.interval.Start("Close").End()

	err := c.conn.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	c.clearWaiters()

	return nil
}
