package storage

import (
	"bytes"
	"encoding/hex"
	"sync"
	"sync/atomic"

	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/pkg/errors"
)

type RAM struct {
	tokens [][]byte
	roomID uint64
	mu     sync.Mutex
}

// implements interface.
var _ Storage = (*RAM)(nil)

func NewRAM() *RAM {
	return &RAM{}
}

func (s *RAM) addDirect(token []byte) {
	s.mu.Lock()
	s.tokens = append(s.tokens, token)
	s.mu.Unlock()
}

func (s *RAM) Add(token []byte) {
	clone := make([]byte, len(token))
	copy(clone, token)
	s.addDirect(clone)
}

func (s *RAM) AddFromString(token string) {
	s.addDirect([]byte(token))
}

func (s *RAM) AddFromHex(token string) {
	data, err := hex.DecodeString(token)
	if err != nil {
		return
	}

	s.addDirect(data)
}

func (s *RAM) Validate(token []byte) (uint64, error) {
	tokens := s.tokens

	for i, t := range tokens {
		if bytes.Equal(t, token) {
			return uint64(i + 1), nil
		}
	}

	return 0, errors.Wrap(constants.ErrInvalid, "token")
}

func (s *RAM) NewRoom() uint64 {
	return atomic.AddUint64(&s.roomID, 1)
}
