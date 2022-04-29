package master

import (
	"log"
	"sync"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/pkg/errors"
)

type SessionStorage struct {
	Clients map[uint64]*SessionServer
	Logger  *log.Logger
	mu      sync.Mutex
}

func (s *SessionStorage) Handle(conn *channel.Channel) {
	server := SessionServer{
		conn:    conn,
		handler: s,
		logger:  s.Logger,
	}

	server.Serve()
}

func (s *SessionStorage) Register(id uint64, client *SessionServer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if prev, ok := s.Clients[id]; ok {
		err := errors.Wrap(client.FlushInstance(prev), "flush error")
		if err != nil {
			s.Logger.Printf("%v\n", err)
		}
	}

	s.Clients[id] = client
}
