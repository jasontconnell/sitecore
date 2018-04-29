package data

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Item struct {
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
}

func (item *Item) GetId() uuid.UUID {
	return item.ID
}

func (item *Item) GetChildren() []ItemNode {
	return item.Children
}

func (item *Item) AddChild(node ItemNode) {
	item.Children = append(item.Children, node)
}

func (item *Item) GetPath() string {
	return item.Path
}

func (item *Item) SetPath(p string) {
	item.Path = p
}

func (item *Item) GetName() string {
	return item.Name
}

func (item *Item) GetParentId() uuid.UUID {
	return item.ParentID
}

func (item *Item) GetParent() ItemNode {
	return item.Parent
}

func (item *Item) SetParent(node ItemNode) {
	item.Parent = item
}

func (item *Item) GetMasterId() uuid.UUID {
	return item.MasterID
}

func (item *Item) GetTemplateId() uuid.UUID {
	return item.TemplateID
}

func (t *Item) String() string {
	return fmt.Sprintf("ID: %v\nName:%v\nPath:%v\n", t.ID, t.Name, t.Path)
}
