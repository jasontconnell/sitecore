package api

import (
	"encoding/binary"

	"github.com/google/uuid"
)

func TryParseUUID(s string) (uuid.UUID, error) {
	u, err := uuid.Parse(trim(s))
	if err != nil {
		return uuid.UUID{}, err
	}
	return u, nil
}

func MustParseUUID(s string) uuid.UUID {
	return uuid.Must(uuid.Parse(trim(s)))
}

func MustParseUUIDProto(lo, hi uint64) uuid.UUID {
	uid, err := ParseUUIDProto(lo, hi)
	if err != nil {
		panic(err)
	}
	return uid
}

func ParseUUIDProto(lo, hi uint64) (uuid.UUID, error) {
	if lo == 0 && hi == 0 {
		return uuid.Nil, nil
	}

	var b uint32 = uint32(lo >> 32)
	var a uint32 = uint32(lo)

	var h uint32 = uint32(hi >> 32)
	var d uint32 = uint32(hi)

	var bytes []byte
	bytes = binary.BigEndian.AppendUint32(bytes, uint32(a))
	bytes = binary.BigEndian.AppendUint16(bytes, uint16(b))
	bytes = binary.BigEndian.AppendUint16(bytes, uint16(b>>16))

	var bsub []byte = make([]byte, 2)
	binary.BigEndian.PutUint16(bsub, uint16(d))
	bytes = append(bytes, bsub[1], bsub[0])

	binary.BigEndian.PutUint16(bsub, uint16(d>>16))
	bytes = append(bytes, bsub[1], bsub[0])

	binary.BigEndian.PutUint16(bsub, uint16(h))
	bytes = append(bytes, bsub[1], bsub[0])

	binary.BigEndian.PutUint16(bsub, uint16(h>>16))
	bytes = append(bytes, bsub[1], bsub[0])

	return uuid.FromBytes(bytes)
}

func trim(s string) string {
	if len(s) < 3 {
		return s
	}
	if s[0] == '{' {
		return string(s[1 : len(s)-1])
	}
	return s
}
