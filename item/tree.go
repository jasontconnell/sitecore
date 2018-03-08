package item

import (
	"github.com/google/uuid"
)

func LoadItemMap(list []*ItemNode) (root *ItemNode, m ItemMap) {
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

func getMap(list []*ItemNode) (root *ItemNode, m ItemMap) {
	m = make(map[uuid.UUID]*ItemNode, len(list))

	rootID, uiderr := uuid.Parse("00000000-0000-0000-0000-000000000000")
	if uiderr != nil {
		return nil, nil
	}

	for _, item := range list {
		m[item.ID] = item
	}

	root = nil
	for _, item := range m {
		if p, ok := m[item.ParentID]; ok {
			p.Children = append(p.Children, item)
			item.Parent = p
		} else if item.ParentID == rootID {
			root = item
		}
	}
	return root, m
}

func setPaths(root *ItemNode) {
	if root.Parent == nil {
		root.Path = "/" + root.Name
	}

	for _, item := range root.Children {
		item.Path = root.Path + "/" + item.Name
		setPaths(item)
	}
}
