package proto

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// Room info for clients connections.
type Room struct {
	ID       uint64          `json:"id"`
	Clients  []*Client       `json:"clients"`
	Endpoint string          `json:"endpoint,omitempty"`
	Result   json.RawMessage `json:"result,omitempty"`
	Error    string          `json:"error,omitempty"`
	ServerID uint64          `json:"-"`
}

func (r *Room) Payload() []byte {
	res, _ := json.Marshal(r)

	return res
}

func (r *Room) Read(data []byte) error {
	return errors.WithStack(json.Unmarshal(data, r))
}
