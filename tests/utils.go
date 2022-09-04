package tests

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) context.Context {
	ctx, cancelFn := context.WithCancel(context.Background())

	t.Cleanup(cancelFn)

	return ctx
}
