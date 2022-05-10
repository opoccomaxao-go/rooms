package proto

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Auth struct {
	Version string `json:"version"`
	Token   string `json:"token"`
}

func (a *Auth) Payload() []byte {
	res, _ := json.Marshal(a)

	return res
}

func (a *Auth) Read(data []byte) error {
	return errors.WithStack(json.Unmarshal(data, a))
}
