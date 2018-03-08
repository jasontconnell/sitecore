package item

import (
	"github.com/google/uuid"
)

type FieldValue struct {
	FieldValueID uuid.UUID
	ItemID       uuid.UUID
	FieldID      uuid.UUID
	Name         string
	Value        string
	Language     string
	Version      int64
	Source       string
}

type FieldValueMap map[uuid.UUID][]*FieldValue
