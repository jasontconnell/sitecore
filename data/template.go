package data

import (
	"github.com/google/uuid"
)

type TemplateQueryRow struct {
	ID              uuid.UUID
	Name            string
	TemplateID      uuid.UUID
	ParentID        uuid.UUID
	MasterID        uuid.UUID
	BaseTemplateIds []uuid.UUID
	Type            string
	Shared          string
	Unversioned     string

	Path     string
	Children []*TemplateQueryRow
}

type template struct {
	ID            uuid.UUID
	Name          string
	Path          string
	Fields        []TemplateFieldNode
	BaseTemplates []TemplateNode
}

func NewTemplateNode(id uuid.UUID, name string, path string) TemplateNode {
	template := &template{ID: id, Name: name, Path: path}
	return template
}

func (t template) GetId() uuid.UUID {
	return t.ID
}

func (t template) GetName() string {
	return t.Name
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

type TemplateNode interface {
	GetId() uuid.UUID
	GetName() string
	AddField(tf TemplateFieldNode)
	GetFields() []TemplateFieldNode
	GetBaseTemplates() []TemplateNode
	AddBaseTemplate(base TemplateNode)
}

type TemplateFieldNode interface {
	GetId() uuid.UUID
	GetName() string
	GetType() string
	GetSource() string
}

func NewTemplateField(id uuid.UUID, name, t, source string) TemplateFieldNode {
	tf := &templateField{ID: id, Name: name, Type: t, Source: source}
	return tf
}

type templateField struct {
	ID     uuid.UUID
	Name   string
	Type   string
	Source string
}

func (tf templateField) GetId() uuid.UUID {
	return tf.ID
}

func (tf templateField) GetName() string {
	return tf.Name
}

func (tf templateField) GetType() string {
	return tf.Type
}

func (tf templateField) GetSource() string {
	return tf.Source
}
