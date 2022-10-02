package engtest

import (
	"github.com/opoccomaxao-go/rooms/engine"
)

func New() engine.Factory {
	return &Factory{}
}

var _ engine.Factory = (*Factory)(nil)

type Factory struct{}

func (f *Factory) New() engine.Engine {
	return &Engine{}
}
