package api

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jasontconnell/sitecore/data"
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
		p, ok := trmap[tr.ParentID]
		if !ok {
			if tr.ParentID == RootID {
				root = tr
				root.Path = "/sitecore"
				continue
			}
			return nil, fmt.Errorf("ParentID not found in map, %v", tr.ParentID)
		}
		p.Children = append(p.Children, tr)
	}

	if root == nil {
		return nil, fmt.Errorf("No root found")
	}

	setTemplatePaths(root)

	templates := loadTemplateData(trmap)

	tnmap := make(map[uuid.UUID]data.TemplateNode)
	for _, t := range templates {
		tnmap[t.GetId()] = t
	}

	for _, t := range tnmap {
		tr, _ := trmap[t.GetId()]
		mapBaseTemplates(tnmap, t, tr)
	}

	return templates, nil
}

func setTemplatePaths(root *data.TemplateQueryRow) {
	for _, c := range root.Children {
		c.Path = root.Path + "/" + root.Name

		setTemplatePaths(c)
	}
}

func loadTemplateData(m map[uuid.UUID]*data.TemplateQueryRow) []data.TemplateNode {
	templates := []data.TemplateNode{}
	for _, tmp := range m {
		if tmp.TemplateID == TemplateID {
			tn := data.NewTemplateNode(tmp.ID, tmp.Name, tmp.Path)
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
		if c.TemplateID == TemplateSectionID {
			flds = append(flds, getFields(c, c.Children)...)
		} else if c.TemplateID == TemplateFieldID {
			s := "Versioned"
			if c.Shared == "1" {
				s = "Shared"
			} else if c.Unversioned == "1" {
				s = "Unversioned"
			}
			tf := data.NewTemplateField(c.ID, c.Name, c.Type, s)
			flds = append(flds, tf)
		}
	}

	return flds
}

func mapBaseTemplates(m map[uuid.UUID]data.TemplateNode, tmp data.TemplateNode, trow *data.TemplateQueryRow) {
	if len(trow.BaseTemplateIds) == 0 {
		return
	}

	for _, id := range trow.BaseTemplateIds {
		if t, ok := m[id]; ok {
			tmp.AddBaseTemplate(t)
		}
	}
}
