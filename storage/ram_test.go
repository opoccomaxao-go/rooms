package storage

import (
	"math"
	"sync/atomic"
	"testing"

	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/stretchr/testify/require"
)

func TestRAM_Validate(t *testing.T) {
	t.Parallel()

	storage := NewRAM()

	tokensS := []string{
		"",
		"test",
		"asadasdgfjsgahsjkfdjksfdds",
	}

	tokensB := [][]byte{
		{1, 2, 3, 4},
		make([]byte, 10000),
	}

	tokensHex := []string{
		"0102030405",
	}

	checkValid := [][]byte{
		{},
		{1, 2, 3, 4, 5},
	}

	checkInvalid := [][]byte{
		{0},
		make([]byte, 10001),
	}

	for _, token := range tokensS {
		storage.AddFromString(token)
	}

	for _, token := range tokensB {
		storage.Add(token)
	}

	for _, token := range tokensHex {
		storage.AddFromHex(token)
	}

	for _, token := range tokensS {
		id, err := storage.Validate([]byte(token))
		require.NoError(t, err)
		require.NotZero(t, id)
	}

	for _, token := range checkValid {
		id, err := storage.Validate(token)
		require.NoError(t, err)
		require.NotZero(t, id)
	}

	for _, token := range checkInvalid {
		id, err := storage.Validate(token)
		require.ErrorIs(t, err, constants.ErrInvalid)
		require.Zero(t, id)
	}
}

func TestRAM_NewRoom(t *testing.T) {
	t.Parallel()

	storage := NewRAM()

	const TotalIDs = math.MaxUint16 + 1

	var res [TotalIDs]int64

	for i := 0; i < TotalIDs; i++ {
		go func() {
			id := storage.NewRoom()
			if id < TotalIDs {
				atomic.AddInt64(&res[id], 1)
			}
		}()
	}

	for i, v := range res {
		require.GreaterOrEqual(t, int64(1), v, i)
	}
}
