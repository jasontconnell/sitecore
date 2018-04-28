package item

import (
	"github.com/google/uuid"
)

func LoadItemMap(list []ItemNode) (root ItemNode, m ItemMap) {
	root, m = getMap(list)
	if root != nil {
		setPaths(root)
	}
	return root, m
}

func LoadFieldMap(list []*FieldValue) FieldValueMap {
	m := FieldValueMap{}
	for _, fv := range list {
		l, ok := m[fv.ItemID]
		if !ok {
			l = []*FieldValue{}
		}
		l = append(l, fv)
		m[fv.ItemID] = l
	}
	return m
}

func getMap(list []ItemNode) (root ItemNode, m ItemMap) {
	m = make(map[uuid.UUID]ItemNode, len(list))

	rootID, uiderr := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if uiderr != nil {
		return nil, nil
	}

	for _, item := range list {
		id := item.GetId()
		m[id] = item
	}

	root = nil
	for _, item := range m {
		if p, ok := m[item.GetParentId()]; ok {
			p.AddChild(item)
			item.SetParent(p)
		} else if item.GetParentId() == rootID {
			root = item
		}
	}
	return root, m
}

func setPaths(root ItemNode) {
	if root.GetParent() == nil {
		root.SetPath("/" + root.GetName())
	}

	for _, item := range root.GetChildren() {
		item.SetPath(root.GetPath() + "/" + item.GetName())
		setPaths(item)
	}
}
