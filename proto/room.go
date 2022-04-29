package proto

// Room info for clients connections.
type Room struct {
	ID      uint64
	Clients []*Client
}
