package proto

import "encoding/json"

// Room info for clients connections.
type Room struct {
	ID       uint64          `json:"id"`
	Clients  []*Client       `json:"clients"`
	Endpoint string          `json:"endpoint,omitempty"`
	Result   json.RawMessage `json:"result,omitempty"`
	ServerID uint64          `json:"-"`
}

func (r *Room) Payload() []byte {
	res, _ := json.Marshal(r)

	return res
}
