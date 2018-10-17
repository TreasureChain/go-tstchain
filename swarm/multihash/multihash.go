package multihash

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	defaultMultihashLength   = 32
	defaultMultihashTypeCode = 0x1b
)

var (
	multihashTypeCode uint8
	MultihashLength   = defaultMultihashLength
)

func init() {
	multihashTypeCode = defaultMultihashTypeCode
	MultihashLength = defaultMultihashLength
}

// check if valid swarm multihash
func isSwarmMultihashType(code uint8) bool {
	return code == multihashTypeCode
}

// GetMultihashLength returns the digest length of the provided multihash
// It will fail if the multihash is not a valid swarm mulithash
func GetMultihashLength(data []byte) (int, int, error) {
	cursor := 0
	typ, c := binary.Uvarint(data)
	if c <= 0 {
		return 0, 0, errors.New("unreadable hashtype field")
	}
	if !isSwarmMultihashType(uint8(typ)) {
		return 0, 0, fmt.Errorf("hash code %x is not a swarm hashtype", typ)
	}
	cursor += c
	hashlength, c := binary.Uvarint(data[cursor:])
	if c <= 0 {
		return 0, 0, errors.New("unreadable length field")
	}
	cursor += c

	// we cheekily assume hashlength < maxint
	inthashlength := int(hashlength)
	if len(data[c:]) < inthashlength {
		return 0, 0, errors.New("length mismatch")
	}
	return inthashlength, cursor, nil
}

// FromMulithash returns the digest portion of the multihash
// It will fail if the multihash is not a valid swarm multihash
func FromMultihash(data []byte) ([]byte, error) {
	hashLength, _, err := GetMultihashLength(data)
	if err != nil {
		return nil, err
	}
	return data[len(data)-hashLength:], nil
}

// ToMulithash wraps the provided digest data with a swarm mulithash header
func ToMultihash(hashData []byte) []byte {
	buf := bytes.NewBuffer(nil)
	b := make([]byte, 8)
	c := binary.PutUvarint(b, uint64(multihashTypeCode))
	buf.Write(b[:c])
	c = binary.PutUvarint(b, uint64(len(hashData)))
	buf.Write(b[:c])
	buf.Write(hashData)
	return buf.Bytes()
}
