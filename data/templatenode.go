package data

import "github.com/google/uuid"

type TemplateNode interface {
	GetId() uuid.UUID
	GetName() string
	GetPath() string

	AddField(tf TemplateFieldNode)
	GetFields() []TemplateFieldNode
	GetAllFields() []TemplateFieldNode
	GetField(id uuid.UUID) TemplateFieldNode

	GetBaseTemplates() []TemplateNode
	AddBaseTemplate(base TemplateNode)
	InheritsTemplate(id uuid.UUID) bool

	HasStandardValues() bool
	GetStandardValuesId() uuid.UUID

	SetStandardValues(item ItemNode)
	GetStandardValues() ItemNode
}

type TemplateFieldNode interface {
	GetId() uuid.UUID
	GetName() string
	GetType() string
	GetSource() FieldSource
}

type TemplateMap map[uuid.UUID]TemplateNode
