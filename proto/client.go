package proto

import "encoding/json"

// Client any active entity.
type Client struct {
	ID    uint64 `json:"id"`
	Token []byte `json:"token,omitempty"`
}

func (c *Client) Payload() []byte {
	res, _ := json.Marshal(c)

	return res
}
