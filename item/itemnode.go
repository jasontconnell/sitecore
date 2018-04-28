package item

import (
	"github.com/google/uuid"
)

type ItemNode interface {
	GetParentId() uuid.UUID
	GetId() uuid.UUID
	GetTemplateId() uuid.UUID

	GetChildren() []ItemNode
	AddChild(node ItemNode)

	GetParent() ItemNode
	SetParent(node ItemNode)

	GetPath() string
	SetPath(p string)

	GetName() string
}

type ItemMap map[uuid.UUID]ItemNode

func (m ItemMap) FindItems(name string) []ItemNode {
	items := []ItemNode{}
	for _, item := range m {
		if item.GetName() == name {
			items = append(items, item)
		}
	}
	return items
}

func (m ItemMap) FindItemByPath(path string) ItemNode {
	var node ItemNode
	for _, item := range m {
		if item.GetPath() == path {
			node = item
		}
	}
	return node
}

func (m ItemMap) FindItemsByTemplate(id uuid.UUID) []ItemNode {
	list := []ItemNode{}
	for _, item := range m {
		if item.GetTemplateId() == id {
			list = append(list, item)
		}
	}

	return list
}
