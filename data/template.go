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
	Namespace     string
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

func (t Template) GetNamespace() string {
	return t.Namespace
}

func (t *Template) SetNamespace(ns string) {
	t.Namespace = ns
}

type TemplateNode interface {
	ItemNode
	GetFields() []TemplateFieldNode
	GetBaseTemplates() []TemplateNode
	GetBaseTemplateIds() []uuid.UUID
	GetNamespace() string
	SetNamespace(ns string)
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
