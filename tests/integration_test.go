package tests

import (
	"context"
	"testing"
	"time"

	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/opoccomaxao-go/rooms/engine/engtest"
	"github.com/opoccomaxao-go/rooms/master"
	"github.com/opoccomaxao-go/rooms/proto"
	"github.com/opoccomaxao-go/rooms/session"
	"github.com/opoccomaxao-go/rooms/storage"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlow(t *testing.T) {
	t.Parallel()

	ctx := TestContext(t)

	ctx, cancelFn := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFn()

	const (
		AuthToken = "12345"
		UserID    = 1
	)

	logger := zerolog.New(zerolog.NewConsoleWriter(
		func(w *zerolog.ConsoleWriter) {
			w.NoColor = true
			w.FormatTimestamp = func(i interface{}) string {
				return time.Now().Format(time.Stamp)
			}
		},
	)).Level(zerolog.TraceLevel)

	storage := storage.NewRAM()
	storage.Add(AuthToken)
	storage.SetVersion(constants.Version)

	mainServer, err := master.New(master.Config{
		Logger:         &logger,
		Storage:        storage,
		SessionAddress: constants.DefaultAddress,
		CreateTimeout:  15 * time.Second,
	})
	require.NoError(t, err)

	go func() {
		require.NoError(t, mainServer.Serve(ctx))
	}()

	time.Sleep(time.Second) // wait for main

	sessionServer, err := session.New(session.Config{
		MasterAddress: constants.DefaultAddress,
		Token:         []byte(AuthToken),
		EngineFactory: engtest.New(),
		Logger:        &logger,
	})
	require.NoError(t, err)

	go func() {
		require.NoError(t, sessionServer.Serve(ctx))
	}()

	time.Sleep(time.Second) // wait for session

	room, err := mainServer.CreateRoom(ctx, []uint64{UserID})
	require.NoError(t, err)
	require.NotNil(t, room)

	assert.NotZero(t, room.ID)
	assert.Equal(t, &proto.Room{
		ID: room.ID,
		Clients: []*proto.Client{
			{ID: UserID},
		},
		ServerID: 1,
	}, room)

	// TODO: implement.
}
