package item

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
}

type ItemNode struct {
	Item
	Children []*ItemNode
	Parent   *ItemNode
}

type ItemMap map[uuid.UUID]*ItemNode

func (t ItemNode) String() string {
	return fmt.Sprintf("ID: %v\nName:%v\nPath:%v\n", t.ID, t.Name, t.Path)
}

func (m ItemMap) FindItems(name string) []*ItemNode {
	items := []*ItemNode{}
	for _, item := range m {
		if item.Name == name {
			items = append(items, item)
		}
	}
	return items
}

func (m ItemMap) FindItemByPath(path string) *ItemNode {
	node := &ItemNode{}
	for _, item := range m {
		if item.Path == path {
			node = item
		}
	}
	return node
}

func (m ItemMap) FindItemsByTemplate(id uuid.UUID) []*ItemNode {
	list := []*ItemNode{}
	for _, item := range m {
		if item.TemplateID == id {
			list = append(list, item)
		}
	}

	return list
}
