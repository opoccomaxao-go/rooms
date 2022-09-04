package session

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/rooms/apm"
	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/opoccomaxao-go/rooms/proto"
	"github.com/opoccomaxao-go/rooms/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Server struct {
	config     Config
	interval   apm.DebuggableInterval
	masterConn *connWrapper
	rooms      []*roomWrapper

	condRooms *sync.Cond

	mu sync.RWMutex
}

type Config struct {
	MasterAddress    string        // MasterAddress is address of master.Server
	Token            []byte        // Token is auth token.
	ReconnectTimeout time.Duration // optional. Default = constants.DefaultTimeoutReconnect

	Logger *zerolog.Logger
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

	if cfg.Logger == nil {
		logger := zerolog.Nop()
		cfg.Logger = &logger
	}

	res := &Server{
		config:     cfg,
		interval:   apm.NewZerologInterval(cfg.Logger, "session.Server."),
		masterConn: &connWrapper{},
		condRooms:  sync.NewCond(&sync.Mutex{}),
	}

	res.masterConn.parent = res
	res.masterConn.init()

	channel, err := channel.Dial(channel.ClientConfig{
		Handler:   res.masterConn.Handler(),
		Address:   cfg.MasterAddress,
		Reconnect: true,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res.masterConn.conn = channel

	return res, nil
}

func (s *Server) Serve(ctx context.Context) error {
	defer s.interval.Start("Serve").End()

	utils.WithContext(ctx).
		AsyncOnDone(func() {
			err := s.masterConn.Close()

			s.config.Logger.Err(err).Stack().Send()
		})

	return s.masterConn.Serve()
}

func (s *Server) Close() error {
	defer s.interval.Start("Close").End()

	return s.masterConn.Close()
}

func (s *Server) AuthClient(token []byte, client net.Conn) error {
	defer s.interval.Start("AuthClient").End()

	// TODO: check token

	// TODO: add to room

	return nil
}

func (s *Server) getCapacity() uint64 {
	defer s.interval.Start("getCapacity").End()

	// TODO: add normal calculations

	return 1
}

func (s *Server) onAuthError(errText string) {
	defer s.interval.Start("onAuthError").End()

	s.config.Logger.Err(errors.New(errText)).Send()

	err := s.Close()
	if err != nil {
		s.config.Logger.Err(err).Stack().Send()
	}
}

func (s *Server) onRoomCreate(room *proto.Room) {
	defer s.interval.Start("onRoomCreate").End()

	// TODO: implement
}

func (s *Server) onRoomCancel(roomID uint64) {
	defer s.interval.Start("onRoomCancel").End()

	// TODO: implement
}
