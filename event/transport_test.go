package event

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransportRW(t *testing.T) {
	t.Parallel()

	batches := [][]*Event{
		{
			{
				Type:    1,
				Payload: []byte("test"),
			},
		},
		{
			{
				Type:    2,
				Payload: []byte("test 2"),
			},
			{
				Type:    3,
				Payload: []byte("test 3"),
			},
		},
		{
			{
				Type:    4,
				Payload: []byte("test 4"),
			},
			{
				Type:    5,
				Payload: []byte("test 5"),
			},
			{
				Type:    6,
				Payload: []byte("test 6"),
			},
		},
	}

	rwc := bytes.Buffer{}

	rwTransport := NewTransport(&rwc)

	for batchID, batch := range batches {
		for _, event := range batch {
			err := rwTransport.Write(event)
			require.NoError(t, err, batchID)
		}

		res := make([]*Event, len(batch))
		for i := range res {
			event, err := rwTransport.Read()
			require.NoError(t, err, i)

			res[i] = event
		}

		require.Equal(t, batch, res, batchID)
	}
}
