package event

import (
	"bytes"
	"io"
	"testing"

	"github.com/opoccomaxao-go/rooms/mock"
	"github.com/stretchr/testify/require"
)

//nolint:gochecknoglobals
var testEvents = []struct {
	desc  string
	event Event
	bytes []byte
}{
	{
		desc: "empty",
		event: Event{
			Type:    0,
			Payload: []byte{},
		},
		bytes: []byte{0, 0, 0, 0},
	},
	{
		desc: "short",
		event: Event{
			Type:    1,
			Payload: []byte("test"),
		},
		bytes: []byte{1, 0, 4, 0, 't', 'e', 's', 't'},
	},
	{
		desc: "long",
		event: Event{
			Type:    0x2345,
			Payload: make([]byte, 0x1234),
		},
		bytes: append([]byte{0x45, 0x23, 0x34, 0x12}, make([]byte, 0x1234)...),
	},
}

func TestEvent_WriteBinary(t *testing.T) {
	t.Parallel()

	for _, tC := range testEvents {
		tC := tC

		var res bytes.Buffer

		res.Grow(0xffff)

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			require.NoError(t, tC.event.WriteBinary(&res))
			require.Equal(t, tC.bytes, res.Bytes())
		})
	}
}

func TestEvent_ReadBinary(t *testing.T) {
	t.Parallel()

	for _, tC := range testEvents {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			var res Event

			err := res.ReadBinary(bytes.NewReader(tC.bytes))
			if err != nil {
				require.ErrorIs(t, err, io.EOF)
			}

			require.Equal(t, tC.event, res)
		})
	}
}

func BenchmarkEvent_WriteBinary(b *testing.B) {
	for _, tC := range testEvents {
		tC := tC

		b.Run(tC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = tC.event.WriteBinary(io.Discard)
			}
		})
	}
}

func BenchmarkEvent_ReadBinary(b *testing.B) {
	for _, tC := range testEvents {
		tC := tC
		event := Event{Payload: []byte{}}
		reader := mock.NewReader(tC.bytes)

		b.Run(tC.desc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = event.ReadBinary(reader)
			}
		})
	}
}
