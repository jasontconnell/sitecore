package data

import (
	"github.com/google/uuid"
)

type fieldValue struct {
	FieldValueID uuid.UUID
	ItemID       uuid.UUID
	FieldID      uuid.UUID
	Name         string
	Value        string
	Language     string
	Version      int64
	Source       string
}

func NewFieldValue(fieldId, itemId uuid.UUID, name, value, language string, version int64, source string) FieldValueNode {
	fv := &fieldValue{}
	fv.ItemID = itemId
	fv.FieldID = fieldId
	fv.Name = name
	fv.Value = value
	fv.Language = language
	fv.Version = version
	fv.Source = source
	return fv
}

func (fv *fieldValue) GetFieldValueId() uuid.UUID {
	return fv.FieldValueID
}

func (fv *fieldValue) GetItemId() uuid.UUID {
	return fv.ItemID
}

func (fv *fieldValue) GetFieldId() uuid.UUID {
	return fv.FieldID
}

func (fv *fieldValue) GetName() string {
	return fv.Name
}

func (fv *fieldValue) GetValue() string {
	return fv.Value
}

func (fv *fieldValue) GetLanguage() string {
	return fv.Language
}

func (fv *fieldValue) GetVersion() int64 {
	return fv.Version
}

func (fv *fieldValue) GetSource() string {
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
