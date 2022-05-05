package proto

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Stats struct {
	Capacity uint64 `json:"capacity"`
}

func (s *Stats) Payload() []byte {
	res, _ := json.Marshal(s)

	return res
}

func (s *Stats) Read(data []byte) error {
	return errors.WithStack(json.Unmarshal(data, s))
}
