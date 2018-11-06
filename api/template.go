package api

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jasontconnell/sitecore/data"
	"strings"
)

func LoadTemplates(connstr string) ([]data.TemplateNode, error) {
	list, err := loadTemplatesFromDb(connstr)
	if err != nil {
		return nil, err
	}

	trmap := make(map[uuid.UUID]*data.TemplateQueryRow)

	for _, tr := range list {
		trmap[tr.ID] = tr
	}

	var root *data.TemplateQueryRow
	for _, tr := range trmap {
		if tr.ParentID == data.RootID {
			root = tr
		}

		p, ok := trmap[tr.ParentID]
		if ok {
			p.Children = append(p.Children, tr)
		}
	}

	if root == nil {
		return nil, fmt.Errorf("No root found")
	}

	root.Path = "/" + root.Name
	setTemplatePaths(root)

	templates := loadTemplateData(trmap)

	tnmap := make(data.TemplateMap)
	for _, t := range templates {
		tnmap[t.GetId()] = t
	}

	for _, t := range tnmap {
		tr, _ := trmap[t.GetId()]
		mapBaseTemplates(tnmap, t, tr)
	}

	return templates, nil
}

func GetTemplateMap(tlist []data.TemplateNode) data.TemplateMap {
	m := make(data.TemplateMap, len(tlist))
	for _, t := range tlist {
		m[t.GetId()] = t
	}
	return m
}

func SetStandardValues(itemMap data.ItemMap, tmap data.TemplateMap) {
	for _, t := range tmap {
		sv, ok := itemMap[t.GetStandardValuesId()]

		if ok {
			t.SetStandardValues(sv)
		}
	}
}

func FilterTemplateMap(tmap data.TemplateMap, paths []string) data.TemplateMap {
	m := make(data.TemplateMap)
	for _, t := range tmap {
		include := false
		for _, b := range paths {
			negate := b[0] == '-'
			b := strings.TrimPrefix(b, "-")
			if len(b) == 0 {
				continue
			}

			if !include && strings.HasPrefix(t.GetPath(), b) {
				include = !negate
			}
		}

		if include {
			m[t.GetId()] = t
		}
	}
	return m
}

func setTemplatePaths(root *data.TemplateQueryRow) {
	for _, c := range root.Children {
		c.Path = root.Path + "/" + c.Name

		setTemplatePaths(c)
	}
}

func loadTemplateData(m map[uuid.UUID]*data.TemplateQueryRow) []data.TemplateNode {
	templates := []data.TemplateNode{}
	for _, tmp := range m {
		if tmp.TemplateID == data.TemplateID {
			tn := data.NewTemplateNode(tmp.ID, tmp.Name, tmp.Path, tmp.StandardValuesId)
			flds := mapFields(tmp)

			for _, f := range flds {
				tn.AddField(f)
			}

			templates = append(templates, tn)
		}
	}
	return templates
}

func mapFields(tmp *data.TemplateQueryRow) []data.TemplateFieldNode {
	return getFields(tmp, tmp.Children)
}

func getFields(tmp *data.TemplateQueryRow, children []*data.TemplateQueryRow) []data.TemplateFieldNode {
	flds := []data.TemplateFieldNode{}
	for _, c := range children {
		if c.TemplateID == data.TemplateSectionID {
			flds = append(flds, getFields(c, c.Children)...)
		} else if c.TemplateID == data.TemplateFieldID {
			s := "VersionedFields"
			if c.Shared == "1" {
				s = "SharedFields"
			} else if c.Unversioned == "1" {
				s = "UnversionedFields"
			}
			tf := data.NewTemplateField(c.ID, c.Name, c.Type, s)
			flds = append(flds, tf)
		}
	}

	return flds
}

func mapBaseTemplates(m data.TemplateMap, tmp data.TemplateNode, trow *data.TemplateQueryRow) {
	hasStdTemplate := false
	for _, id := range trow.BaseTemplateIds {
		if t, ok := m[id]; ok {
			tmp.AddBaseTemplate(t)
		}

		if id == data.StandardTemplateID {
			hasStdTemplate = true
		}
	}

	stdTemplate, stdTempFound := m[data.StandardTemplateID]
	if !hasStdTemplate && stdTempFound {
		tmp.AddBaseTemplate(stdTemplate)
	}
}
