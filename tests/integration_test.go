package tests

import (
	"log"
	"testing"

	"github.com/opoccomaxao-go/rooms/master"
	"github.com/opoccomaxao-go/rooms/session"
	"github.com/opoccomaxao-go/rooms/storage"
	"github.com/stretchr/testify/require"
)

func TestFlow(t *testing.T) {
	t.Parallel()

	storage := storage.NewRAM()
	storage.Add("12345")

	mainServer, err := master.New(master.Config{
		Logger:  log.Default(),
		Storage: storage,
	})
	require.NoError(t, err)

	sessionServer, err := session.New(session.Config{})
	require.NoError(t, err)

	_ = mainServer
	_ = sessionServer
}
