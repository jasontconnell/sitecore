package data

import (
	"github.com/google/uuid"
)

type ItemData interface {
	GetId() uuid.UUID
	SetId(id uuid.UUID)

	GetParentId() uuid.UUID
	SetParentId(id uuid.UUID)

	GetTemplateId() uuid.UUID
	SetTemplateId(id uuid.UUID)

	GetMasterId() uuid.UUID
	SetMasterId(id uuid.UUID)

	GetLevel() int
	SetLevel(level int)
	GetPath() string
	SetPath(p string)

	GetName() string
	SetName(n string)
}

type ItemNode interface {
	ItemData

	GetChildren() []ItemNode
	AddChild(node ItemNode)

	GetParent() ItemNode
	SetParent(node ItemNode)

	GetFieldValues() []FieldValueNode
	AddFieldValue(fv FieldValueNode)
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
