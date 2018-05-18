package data

import (
	"github.com/google/uuid"
)

type templateMeta struct {
	Type            string
	BaseTemplateIds []uuid.UUID
}

type template struct {
	templateMeta
	ItemNode
	Fields        []TemplateFieldNode
	BaseTemplates []TemplateNode
	Namespace     string
}

func NewTemplateNode(item ItemNode, fldType string, baseTemplateIds []uuid.UUID) TemplateNode {
	template := &template{ItemNode: item, templateMeta: templateMeta{Type: fldType, BaseTemplateIds: baseTemplateIds}}
	return template
}

func (t template) GetFields() []TemplateFieldNode {
	return t.Fields
}

func (t *template) AddField(fld TemplateFieldNode) {
	t.Fields = append(t.Fields, fld)
}

func (t *template) AddBaseTemplate(base TemplateNode) {
	t.BaseTemplates = append(t.BaseTemplates, base)
}

func (t template) GetBaseTemplates() []TemplateNode {
	return t.BaseTemplates
}

func (t template) GetBaseTemplateIds() []uuid.UUID {
	return t.BaseTemplateIds
}

func (t template) GetType() string {
	return t.Type
}

type TemplateNode interface {
	ItemNode
	AddField(tf TemplateFieldNode)
	GetFields() []TemplateFieldNode
	GetBaseTemplates() []TemplateNode
	GetBaseTemplateIds() []uuid.UUID
	AddBaseTemplate(base TemplateNode)
	GetType() string
}

type TemplateFieldNode interface {
	ItemNode
	GetType() string
}

func NewTemplateField(inner ItemNode, t string) TemplateFieldNode {
	tf := &templateField{ItemNode: inner, Type: t}
	return tf
}

type templateField struct {
	ItemNode
	Type string
}

func (tf templateField) GetType() string {
	return tf.Type
}
