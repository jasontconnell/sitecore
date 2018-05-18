package data

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

func NewFieldValue(fvId, itemId, fieldId uuid.UUID, name, value, language string, version int64, source string) FieldValueNode {
	fv := &FieldValue{}
	fv.FieldValueID = fvId
	fv.ItemID = itemId
	fv.FieldID = fieldId
	fv.Name = name
	fv.Value = value
	fv.Language = language
	fv.Version = version
	fv.Source = source
	return fv
}

func (fv *FieldValue) GetFieldValueId() uuid.UUID {
	return fv.FieldValueID
}

func (fv *FieldValue) GetItemId() uuid.UUID {
	return fv.ItemID
}

func (fv *FieldValue) GetFieldId() uuid.UUID {
	return fv.FieldID
}

func (fv *FieldValue) GetName() string {
	return fv.Value
}

func (fv *FieldValue) GetValue() string {
	return fv.Value
}

func (fv *FieldValue) GetLanguage() string {
	return fv.Language
}

func (fv *FieldValue) GetVersion() int64 {
	return fv.Version
}

func (fv *FieldValue) GetSource() string {
	return fv.Source
}

type FieldValueMap map[uuid.UUID][]FieldValueNode

type FieldValueNode interface {
	GetFieldValueId() uuid.UUID
	GetItemId() uuid.UUID
	GetFieldId() uuid.UUID
	GetName() string
	GetValue() string
	GetLanguage() string
	GetVersion() int64
	GetSource() string
}
