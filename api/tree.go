package api

import (
	"github.com/google/uuid"
	"github.com/jasontconnell/sitecore/data"
	"path"
	"sort"
	"strings"
)

var rootID uuid.UUID = MustParseUUID("00000000-0000-0000-0000-000000000000")

func LoadItemMap(list []data.ItemNode) (root data.ItemNode, m data.ItemMap) {
	root, m = getMap(list)
	if root != nil {
		setTreeData(root, 0)
	}
	return root, m
}

func LoadFieldMap(list []data.FieldValueNode) data.FieldValueMap {
	m := data.FieldValueMap{}
	for _, fv := range list {
		l, ok := m[fv.GetItemId()]
		if !ok {
			l = []data.FieldValueNode{}
		}
		l = append(l, fv)
		m[fv.GetItemId()] = l
	}
	return m
}

func FilterItemMap(m data.ItemMap, paths []string) data.ItemMap {
	filteredMap := make(data.ItemMap)
	for _, item := range m {
		include := false
		for _, b := range paths {
			negate := b[0] == '-'
			b := strings.TrimPrefix(b, "-")

			if !include && strings.HasPrefix(item.GetPath(), b) {
				include = !negate
				break
			} else {
				parent := path.Dir(b)

				for parent != "/" && parent != "" && !include {
					include = item.GetPath() == parent
					parent = path.Dir(parent)
				}
			}
		}

		if include {
			filteredMap[item.GetId()] = item
		}
	}

	return filteredMap
}

func AssignFieldValues(m data.ItemMap, values []data.FieldValueNode) {
	sort.Slice(values, func(i, j int) bool {
		return values[i].GetName() < values[j].GetName()
	})

	for _, fv := range values {
		if item, ok := m[fv.GetItemId()]; ok {
			item.AddFieldValue(fv)
		}
	}
}

func getMap(list []data.ItemNode) (root data.ItemNode, m data.ItemMap) {
	m = make(map[uuid.UUID]data.ItemNode, len(list))

	for _, item := range list {
		id := item.GetId()
		m[id] = item
	}

	root = nil
	for _, item := range m {
		pid := item.GetParentId()
		if pid == rootID {
			root = item
		} else if p, ok := m[pid]; ok {
			p.AddChild(item)
			item.SetParent(p)
		}
	}
	return root, m
}

func setTreeData(root data.ItemNode, level int) {
	if root.GetParent() == nil {
		root.SetPath("/" + root.GetName())
		root.SetLevel(0)
	}

	for _, item := range root.GetChildren() {
		item.SetPath(root.GetPath() + "/" + item.GetName())
		item.SetLevel(level + 1)
		setTreeData(item, level + 1)
	}
}
