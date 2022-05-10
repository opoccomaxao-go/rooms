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

	const Version = "1.1.1"

	storage := NewRAM()
	storage.SetVersion(Version)

	tokens := []string{
		"",
		"test",
		"asadasdgfjsgahsjkfdjksfdds",
	}

	checkInvalidTokens := []string{
		"\x00",
		string(make([]byte, 10001)),
	}

	checkInvalidVersions := []string{
		"",
		"1.1.1.1",
	}

	for _, token := range tokens {
		storage.Add(token)
	}

	for _, token := range tokens {
		id, err := storage.Validate(Version, token)
		require.NoError(t, err)
		require.NotZero(t, id)
	}

	for _, token := range checkInvalidTokens {
		id, err := storage.Validate(Version, token)
		require.ErrorIs(t, err, constants.ErrInvalid)
		require.Zero(t, id)
	}

	for _, version := range checkInvalidVersions {
		id, err := storage.Validate(version, tokens[0])
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
