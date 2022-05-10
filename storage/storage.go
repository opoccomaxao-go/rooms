package storage

type Storage interface {
	Validate(version string, token string) (uint64, error)
	NewRoom() uint64
}
