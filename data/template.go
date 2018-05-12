package data

import (
	"github.com/google/uuid"
)

type TemplateMeta struct {
	Type            string
	BaseTemplateIds []uuid.UUID
}

type Template struct {
	TemplateMeta
	Item
	Fields        []TemplateFieldNode
	BaseTemplates []TemplateNode
}

func (t Template) GetFields() []TemplateFieldNode {
	return t.Fields
}

func (t Template) GetBaseTemplates() []TemplateNode {
	return t.BaseTemplates
}

func (t Template) GetBaseTemplateIds() []uuid.UUID {
	return t.BaseTemplateIds
}

type TemplateNode interface {
	ItemNode
	GetFields() []TemplateFieldNode
	GetBaseTemplates() []TemplateNode
	GetBaseTemplateIds() []uuid.UUID
}

type TemplateFieldNode interface {
	ItemNode
	GetType() string
}

type TemplateField struct {
	ItemNode
	Type string
}

func (tf TemplateField) GetType() string {
	return tf.Type
}
