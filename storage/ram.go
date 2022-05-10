package storage

import (
	"sync"
	"sync/atomic"

	"github.com/opoccomaxao-go/rooms/constants"
	"github.com/pkg/errors"
)

type RAM struct {
	tokens  []string
	version string
	roomID  uint64
	mu      sync.Mutex
}

// implements interface.
var _ Storage = (*RAM)(nil)

func NewRAM() *RAM {
	return &RAM{}
}

func (s *RAM) addDirect(token string) {
	s.mu.Lock()
	s.tokens = append(s.tokens, token)
	s.mu.Unlock()
}

func (s *RAM) SetVersion(version string) {
	s.version = version
}

func (s *RAM) Add(token string) {
	s.addDirect(token)
}

func (s *RAM) Validate(version string, token string) (uint64, error) {
	if s.version != version {
		return 0, errors.Wrap(constants.ErrInvalid, "version")
	}

	for i, t := range s.tokens {
		if t == token {
			return uint64(i + 1), nil
		}
	}

	return 0, errors.Wrap(constants.ErrInvalid, "token")
}

func (s *RAM) NewRoom() uint64 {
	return atomic.AddUint64(&s.roomID, 1)
}
