package data

import (
	"github.com/google/uuid"
)

type FieldSource string

const (
	VersionedFields   FieldSource = "VersionedFields"
	UnversionedFields FieldSource = "UnversionedFields"
	SharedFields      FieldSource = "SharedFields"
)

type Language string

const (
	None    Language = ""
	English Language = "en"
)

func GetLanguage(lan string) Language {
	if lan == "" {
		return None
	}
	return Language(lan)
}

func (f FieldSource) String() string {
	return string(f)
}

func GetSource(s string) FieldSource {
	return FieldSource(s)
}

type fieldValue struct {
	FieldValueID uuid.UUID
	ItemID       uuid.UUID
	FieldID      uuid.UUID
	Name         string
	Value        string
	Language     Language
	Version      int64
	Source       FieldSource
}

func NewFieldValue(fieldId, itemId uuid.UUID, name, value string, language Language, version int64, source FieldSource) FieldValueNode {
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

func (fv *fieldValue) SetValue(val string) {
	fv.Value = val
}

func (fv *fieldValue) GetLanguage() Language {
	return fv.Language
}

func (fv *fieldValue) GetVersion() int64 {
	return fv.Version
}

func (fv *fieldValue) GetSource() FieldSource {
	return fv.Source
}

type FieldValueMap map[uuid.UUID][]FieldValueNode

type FieldValueNode interface {
	GetFieldValueId() uuid.UUID
	GetItemId() uuid.UUID
	GetFieldId() uuid.UUID
	GetName() string
	GetValue() string
	SetValue(val string)
	GetLanguage() Language
	GetVersion() int64
	GetSource() FieldSource
}
