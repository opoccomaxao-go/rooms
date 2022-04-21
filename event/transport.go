package event

import (
	"io"

	"github.com/pkg/errors"
)

type Transport interface {
	Write(event *Event) error
	Read() (*Event, error)
	Close() error
}

type rw struct {
	stream io.ReadWriter
}

func (rw *rw) Write(event *Event) error {
	return errors.WithStack(event.WriteBinary(rw.stream))
}

func (rw *rw) Read() (*Event, error) {
	var res Event

	err := res.ReadBinary(rw.stream)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &res, nil
}

func (rw *rw) Close() error {
	if c, ok := rw.stream.(io.Closer); ok {
		return errors.WithStack(c.Close())
	}

	return nil
}

func NewTransport(stream io.ReadWriter) Transport {
	return &rw{
		stream: stream,
	}
}
