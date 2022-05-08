package session

import (
	"net"
	"sync"
	"time"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/pkg/errors"
)

type Server struct {
	config     Config
	masterConn *connWrapper
	rooms      []*roomWrapper

	condRooms *sync.Cond

	mu sync.RWMutex
}

type Config struct {
	MasterAddress    string        // MasterAddress is address of master.Server
	Token            []byte        // Token is auth token.
	ReconnectTimeout time.Duration // optional. Default = constants.DefaultTimeoutReconnect
}

func New(cfg Config) (*Server, error) {
	if cfg.MasterAddress == "" {
		cfg.MasterAddress = constants.DefaultAddress
	}

	if len(cfg.Token) == 0 {
		return nil, errors.WithMessage(constants.ErrNoParam, "Token")
	}

	if cfg.ReconnectTimeout <= 0 {
		cfg.ReconnectTimeout = constants.DefaultTimeoutReconnect
	}

	masterConn := &connWrapper{}

	channel, err := channel.Dial(channel.ClientConfig{
		Handler:   masterConn.Handler(),
		Address:   cfg.MasterAddress,
		Reconnect: true,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	masterConn.conn = channel

	return &Server{
		config:     cfg,
		masterConn: masterConn,
		condRooms:  sync.NewCond(&sync.Mutex{}),
	}, nil
}

func (s *Server) AuthClient(token []byte, client net.Conn) error {
	// TODO: check token

	// TODO: add to room

	return nil
}
