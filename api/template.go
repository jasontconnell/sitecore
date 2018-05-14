package api

import (
	"github.com/jasontconnell/sitecore/data"
)

var templateId = MustParseUUID("AB86861A-6030-46C5-B394-E8F99E8B87DB")
var fieldId = MustParseUUID("455A3E98-A627-4B40-8035-E683A0331AC7")
var templateSectionId = MustParseUUID("E269FBB5-3750-427A-9149-7AA950B49301")

func GetTemplates(m data.ItemMap) []data.ItemNode {
	return m.FindItemsByTemplate(templateId)
}

func LoadTemplates(connstr string) ([]data.TemplateNode, error) {
	list, err := loadTemplatesFromDb(connstr)
	if err != nil {
		return nil, err
	}

	itemNodes := []data.ItemNode{}
	for _, t := range list {
		itemNodes = append(itemNodes, t.(data.ItemNode))
	}
	_, m := LoadItemMap(itemNodes) // don't care about root

	LoadTemplateData(m)
	retList := []data.TemplateNode{}

	for _, item := range itemNodes {
		if item.GetTemplateId() == templateId {
			retList = append(retList, item.(data.TemplateNode))
		}
	}

	return retList, nil
}

func LoadTemplateData(m data.ItemMap) {
	for _, tmp := range m {
		t := tmp.(*data.Template)
		if t.TemplateID == templateId {
			getBaseTemplates(m, t)
			mapFields(t)
		}
	}
}

func mapFields(tmp *data.Template) {
	getFields(tmp, tmp.Children)
}

func getFields(tmp *data.Template, children []data.ItemNode) {
	for _, c := range children {
		ct := c.(*data.Template)
		if ct.TemplateID == templateSectionId {
			getFields(tmp, ct.Children)
			return
		} else if ct.TemplateID == fieldId {

			tf := data.TemplateField{Type: ct.Type, ItemNode: c}
			tmp.Fields = append(tmp.Fields, tf)
		}
	}
}

func getBaseTemplates(m data.ItemMap, tmp *data.Template) {
	if len(tmp.BaseTemplateIds) == 0 {
		return
	}

	for _, id := range tmp.BaseTemplateIds {
		if t, ok := m[id]; ok {
			tmp.BaseTemplates = append(tmp.BaseTemplates, t.(*data.Template))
		}
	}
}
