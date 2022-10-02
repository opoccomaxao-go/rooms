package engtest

import (
	"github.com/opoccomaxao-go/rooms/engine"
)

var _ engine.Engine = (*Engine)(nil)

type Engine struct{}
