package master

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/opoccomaxao-go/rooms/proto"
	"github.com/opoccomaxao-go/rooms/storage"
	"github.com/pkg/errors"
)

type Server struct {
	config         Config
	server         *channel.Server
	sessionStorage *SessionStorage

	mu sync.Mutex
}

type Config struct {
	Logger         *log.Logger
	Storage        storage.Storage
	SessionAddress string        // SessionAddress is address for session-server listening.
	CreateTimeout  time.Duration // CreateTimeout is NewRoom timeout.
}

func New(cfg Config) (*Server, error) {
	if cfg.Logger == nil {
		cfg.Logger = log.Default()
	}

	if cfg.Storage == nil {
		cfg.Storage = storage.NewRAM()
	}

	if cfg.SessionAddress == "" {
		cfg.SessionAddress = constants.DefaultAddress
	}

	if cfg.CreateTimeout <= 0 {
		cfg.CreateTimeout = constants.DefaultTimeout
	}

	sessionStorage := newSessionStorage(&cfg)

	server, err := channel.NewServer(channel.ServerConfig{
		Address: cfg.SessionAddress,
		Handler: sessionStorage,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Server{
		config:         cfg,
		sessionStorage: sessionStorage,
		server:         server,
	}, nil
}

func (s *Server) Serve() error {
	return errors.WithStack(s.server.Listen())
}

func (s *Server) NewRoom(userIDs []uint64) (*proto.Room, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	room := &proto.Room{
		ID:      s.config.Storage.NewRoom(),
		Clients: make([]*proto.Client, len(userIDs)),
	}

	for index, id := range userIDs {
		room.Clients[index].ID = id
	}

	ctx, cancelFn := context.WithTimeout(context.TODO(), s.config.CreateTimeout)
	defer cancelFn()

	res, err := s.sessionStorage.CreateRoom(ctx, room)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res, nil
}
