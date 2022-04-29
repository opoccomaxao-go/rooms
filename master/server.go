package master

import (
	"log"
	"sync"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/rooms/constants"
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
	SessionAddress string // SessionAddress is address for session-server listening.
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

	sessionStorage := &SessionStorage{
		Clients: map[uint64]*SessionServer{},
		Logger:  cfg.Logger,
	}

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

func (s *Server) NewRoom(userIDs []uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	panic("implement")
}
