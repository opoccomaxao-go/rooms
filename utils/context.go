package utils

import (
	"context"
)

type Context struct {
	internal context.Context //nolint:containedctx // core feature
}

func WithContext(ctx context.Context) *Context {
	return &Context{
		internal: ctx,
	}
}

func (ctx *Context) OnDone(f func()) *Context {
	<-ctx.internal.Done()

	f()

	return ctx
}

func (ctx *Context) AsyncOnDone(f func()) *Context {
	go ctx.OnDone(f)

	return ctx
}
