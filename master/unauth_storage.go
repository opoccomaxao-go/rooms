package master

type UnauthStorage struct{}

func NewUnauthStorage() Storage {
	return &UnauthStorage{}
}
