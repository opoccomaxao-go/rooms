package storage

type Storage interface {
	Validate(token []byte) (uint64, error)
}
