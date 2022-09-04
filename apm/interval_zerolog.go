package apm

import (
	"github.com/rs/zerolog"
)

var _ DebuggableInterval = (*intervalZerolog)(nil)

type intervalZerolog struct {
	logger *zerolog.Logger
	prefix string
}

func NewZerologInterval(logger *zerolog.Logger, prefix string) DebuggableInterval {
	return &intervalZerolog{
		logger: logger,
		prefix: prefix,
	}
}

func (i *intervalZerolog) Start(name string) Interval {
	logger := i.logger.With().Str("name", i.prefix+name).Logger()
	logger.Debug().Msg("start")

	return func() {
		logger.Debug().Msg("end")
	}
}
