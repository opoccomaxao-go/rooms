package utils

type Channels[T any] struct {
	internal []chan T
}

func WithChannels[T any](channels []chan T) *Channels[T] {
	return &Channels[T]{
		internal: channels,
	}
}

func (c *Channels[T]) Notify(value T) {
	for _, v := range c.internal {
		v <- value
	}
}

func (c *Channels[T]) Close() {
	for _, v := range c.internal {
		close(v)
	}
}
