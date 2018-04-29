package sitecore

import (
    "github.com/google/uuid"
)

func TryParseUUID(s string) (uuid.UUID, error){
    u, err := uuid.Parse(s)
    if err != nil {
        return uuid.UUID{}, err
    }
    return u, nil
}

func MustParseUUID(s string) uuid.UUID {
    return uuid.Must(uuid.Parse(s))
}