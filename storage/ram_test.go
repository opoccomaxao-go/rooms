package storage

import (
	"testing"

	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/stretchr/testify/require"
)

func TestRAM(t *testing.T) {
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
