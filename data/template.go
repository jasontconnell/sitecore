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

	StandardValuesId uuid.UUID
	Type             string
	Shared           string
	Unversioned      string

	Path     string
	Children []*TemplateQueryRow
}

type template struct {
	ID          uuid.UUID
	Name        string
	Path        string
	Fields      []TemplateFieldNode
	AllFieldMap map[uuid.UUID]TemplateFieldNode

	BaseTemplates   []TemplateNode
	BaseTemplateMap map[uuid.UUID]TemplateNode

	StandardValuesID uuid.UUID
	StandardValues   ItemNode
}

func NewTemplateNode(id uuid.UUID, name string, path string, standardValuesId uuid.UUID) TemplateNode {
	template := &template{ID: id, Name: name, Path: path, StandardValuesID: standardValuesId}
	template.BaseTemplateMap = make(map[uuid.UUID]TemplateNode)
	return template
}

func (t template) GetId() uuid.UUID {
	return t.ID
}

func (t template) GetName() string {
	return t.Name
}

func (t template) GetPath() string {
	return t.Path
}

func (t template) GetFields() []TemplateFieldNode {
	return t.Fields
}

func (t *template) GetField(id uuid.UUID) TemplateFieldNode {
	if t.AllFieldMap == nil {
		t.AllFieldMap = make(map[uuid.UUID]TemplateFieldNode)
		list := t.GetAllFields()
		for _, f := range list {
			t.AllFieldMap[f.GetId()] = f
		}
	}

	fld, _ := t.AllFieldMap[id]
	return fld
}

func (t *template) FindField(name string) TemplateFieldNode {
	flds := t.GetAllFields()
	var ret TemplateFieldNode
	for _, f := range flds {
		if f.GetName() == name {
			ret = f
			break
		}
	}
	return ret
}

func (t template) GetAllFields() []TemplateFieldNode {
	v := make(map[uuid.UUID]bool)
	return internalGetAllFields(&t, v)
}

func internalGetAllFields(t TemplateNode, v map[uuid.UUID]bool) []TemplateFieldNode {
	if _, ok := v[t.GetId()]; ok {
		return nil
	}
	list := []TemplateFieldNode{}
	for _, f := range t.GetFields() {
		if _, ok := v[f.GetId()]; ok {
			continue
		}
		list = append(list, f)
	}

	v[t.GetId()] = true

	for _, b := range t.GetBaseTemplates() {
		bf := internalGetAllFields(b, v)
		if bf == nil {
			continue
		}
		list = append(list, bf...)
	}
	return list
}

func (t *template) AddField(fld TemplateFieldNode) {
	t.Fields = append(t.Fields, fld)
}

func (t *template) AddBaseTemplate(base TemplateNode) {
	if _, ok := t.BaseTemplateMap[base.GetId()]; !ok {
		t.BaseTemplates = append(t.BaseTemplates, base)
		t.BaseTemplateMap[base.GetId()] = base
	}
}

func (t template) GetBaseTemplates() []TemplateNode {
	return t.BaseTemplates
}

func (t template) InheritsTemplate(id uuid.UUID) bool {
	vmap := make(map[uuid.UUID]bool)
	return internalInheritsTemplate(&t, id, vmap)
}

func internalInheritsTemplate(t TemplateNode, id uuid.UUID, v map[uuid.UUID]bool) bool {
	if inherits, ok := v[id]; ok {
		return inherits
	}

	for _, b := range t.GetBaseTemplates() {
		if b.GetId() == id {
			v[id] = true
			break
		} else {
			return internalInheritsTemplate(t, id, v)
		}
	}
	return v[id]
}

func (t template) HasStandardValues() bool {
	return t.StandardValuesID != EmptyID
}

func (t template) GetStandardValuesId() uuid.UUID {
	return t.StandardValuesID
}

func (t template) SetStandardValues(item ItemNode) {
	t.StandardValues = item
}

func (t template) GetStandardValues() ItemNode {
	return t.StandardValues
}

func NewTemplateField(id uuid.UUID, name, t string, source FieldSource) TemplateFieldNode {
	tf := &templateField{ID: id, Name: name, Type: t, Source: source}
	return tf
}

type templateField struct {
	ID     uuid.UUID
	Name   string
	Type   string
	Source FieldSource
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

func (tf templateField) GetSource() FieldSource {
	return tf.Source
}
