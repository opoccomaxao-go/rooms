package session

import (
	"time"

	"github.com/opoccomaxao-go/rooms/apm"
	"github.com/opoccomaxao-go/rooms/engine"
	"github.com/opoccomaxao-go/rooms/proto"
	"github.com/rs/zerolog"
)

type roomWrapper struct {
	roomData *proto.Room
	parent   *Server

	logger   zerolog.Logger
	interval apm.DebuggableInterval

	id      uint64
	clients []*clientWrapper
	mapping map[uint64]*clientWrapper
}

func (r *roomWrapper) init() {
	r.logger = r.parent.config.Logger.With().Logger()
	r.interval = apm.NewZerologInterval(&r.logger, "session.roomWrapper.")

	clientsTotal := len(r.roomData.Clients)

	r.id = r.roomData.ID
	r.clients = make([]*clientWrapper, clientsTotal)
	r.mapping = make(map[uint64]*clientWrapper, clientsTotal)

	for i, clientData := range r.roomData.Clients {
		client := clientWrapper{}

		r.clients[i] = &client
		r.mapping[clientData.ID] = &client
	}
}

func (r *roomWrapper) Serve(engine engine.Engine) {
	defer r.interval.Start("Serve").End()
	defer r.finish()

	// TODO: implement.
	time.Sleep(time.Second)
}

func (r *roomWrapper) finish() {
	defer r.interval.Start("finish").End()

	r.parent.onSessionEnd(r.id)
}
