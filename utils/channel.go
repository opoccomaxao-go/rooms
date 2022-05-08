package utils

import "context"

// Channel usage:
//  WithChannel(someChan).
//    OnBeforeClose(func(){ ... }).
//    OnAfterClose(func(){ ... }).
//    AsyncCloseOnDone(ctx)
type Channel[T any] struct {
	internal    chan T
	beforeClose func()
	afterClose  func()
}

func WithChannel[T any](channel chan T) *Channel[T] {
	return &Channel[T]{
		internal: channel,
	}
}

func (c *Channel[T]) OnBeforeClose(f func()) *Channel[T] {
	c.beforeClose = f

	return c
}

func (c *Channel[T]) OnAfterClose(f func()) *Channel[T] {
	c.afterClose = f

	return c
}

func (c *Channel[T]) CloseOnDone(ctx context.Context) {
	<-ctx.Done()
	TryExec(c.beforeClose)
	close(c.internal)
	TryExec(c.afterClose)
}

func (c *Channel[T]) AsyncCloseOnDone(ctx context.Context) {
	go c.CloseOnDone(ctx)
}
