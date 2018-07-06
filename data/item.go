package data

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type item struct {
	ID         uuid.UUID
	Name       string
	TemplateID uuid.UUID
	ParentID   uuid.UUID
	MasterID   uuid.UUID
	Created    time.Time
	Updated    time.Time
	Path       string
	Parent     ItemNode
	Children   []ItemNode
	Level      int

	FieldValues []FieldValueNode
}

func NewBlankItemNode() ItemNode {
	return new(item)
}

func NewItemNode(id uuid.UUID, name string, templateId, parentId, masterId uuid.UUID) ItemNode {
	item := &item{}
	item.ID = id
	item.Name = name
	item.TemplateID = templateId
	item.ParentID = parentId
	item.MasterID = masterId
	return item
}

func (item *item) GetId() uuid.UUID {
	return item.ID
}

func (item *item) SetId(id uuid.UUID) {
	item.ID = id
}

func (item *item) GetChildren() []ItemNode {
	return item.Children
}

func (item *item) AddChild(node ItemNode) {
	item.Children = append(item.Children, node)
}

func (item *item) GetLevel() int {
	return item.Level
}

func (item *item) SetLevel(level int) {
	item.Level = level
}

func (item *item) GetPath() string {
	return item.Path
}

func (item *item) SetPath(p string) {
	item.Path = p
}

func (item *item) GetName() string {
	return item.Name
}

func (item *item) SetName(n string) {
	item.Name = n
}

func (item *item) GetParentId() uuid.UUID {
	return item.ParentID
}

func (item *item) SetParentId(id uuid.UUID) {
	item.ParentID = id
}

func (item *item) GetParent() ItemNode {
	return item.Parent
}

func (item *item) SetParent(node ItemNode) {
	item.Parent = item
}

func (item *item) GetMasterId() uuid.UUID {
	return item.MasterID
}

func (item *item) SetMasterId(id uuid.UUID) {
	item.MasterID = id
}

func (item *item) GetTemplateId() uuid.UUID {
	return item.TemplateID
}

func (item *item) SetTemplateId(id uuid.UUID) {
	item.TemplateID = id
}

func (t *item) String() string {
	return fmt.Sprintf("ID: %v\nName:%v\nPath:%v\n", t.ID, t.Name, t.Path)
}

func (t *item) GetFieldValues() []FieldValueNode {
	return t.FieldValues
}

func (t *item) AddFieldValue(fv FieldValueNode) {
	t.FieldValues = append(t.FieldValues, fv)
}
