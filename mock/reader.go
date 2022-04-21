package mock

import "io"

type Reader struct {
	data []byte
	size int
}

var _ io.Reader = &Reader{}

func NewReader(data []byte) io.Reader {
	return &Reader{
		data: data,
		size: len(data),
	}
}

func (r *Reader) Read(p []byte) (int, error) {
	copy(p, r.data)

	return r.size, nil
}
