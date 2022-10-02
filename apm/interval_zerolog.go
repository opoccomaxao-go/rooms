package apm

import (
	"sync/atomic"

	"github.com/rs/zerolog"
)

var _ DebuggableInterval = (*intervalZerolog)(nil)

//noling:gochecknoglobals // required for debug.
var id uint64

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
	logger := i.logger.With().
		Str("name", i.prefix+name).
		Uint64("id", atomic.AddUint64(&id, 1)).
		Logger()
	logger.Debug().Msg("start")

	return func() {
		logger.Debug().Msg("end")
	}
}
