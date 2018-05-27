package api

import (
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

func trim(s string) string {
	if len(s) < 3 {
		return s
	}
	if s[0] == '{' {
		return string(s[1 : len(s)-1])
	}
	return s
}
