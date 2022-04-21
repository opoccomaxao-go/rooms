package tests

import (
	"testing"

	"github.com/opoccomaxao-go/rooms/master"
	"github.com/opoccomaxao-go/rooms/session"
	"github.com/stretchr/testify/require"
)

func TestFlow(t *testing.T) {
	t.Parallel()

	mainServer, err := master.New(master.Config{})
	require.NoError(t, err)

	sessionServer, err := session.New(session.Config{})
	require.NoError(t, err)

	_ = mainServer
	_ = sessionServer
}
