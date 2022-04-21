package event

import (
	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/pkg/errors"
)

type Server struct {
	config Config
}

func NewServer(config Config) (*Server, error) {
	if config.Transport == nil {
		return nil, errors.WithMessage(constants.ErrNoParam, "Transport")
	}

	return &Server{
		config: config,
	}, nil
}
