package comm

import (
	"encoding/binary"
)

func GetUint32(body []byte) uint32 {
	return binary.BigEndian.Uint32(body)
}

func GetUint64(body []byte) uint64 {
	return binary.BigEndian.Uint64(body)
}

func PutUint32(value uint32, body []byte) {
	binary.BigEndian.PutUint32(body, value)
}

func PutUint64(value uint64, body []byte) {
	binary.BigEndian.PutUint64(body, value)
}
