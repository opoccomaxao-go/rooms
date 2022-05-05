package proto

import (
	"encoding/binary"
)

const uint64Bytes = 8

type ID = uint64

func PayloadID(id ID) []byte {
	res := make([]byte, uint64Bytes)

	binary.BigEndian.PutUint64(res, id)

	return res
}

func ReadID(payload []byte) ID {
	return binary.BigEndian.Uint64(payload)
}
