package item

import (
	"fmt"
)

type Template struct {
	Item          *Item
	Fields        []TemplateField
	BaseTemplates []*Template
}

type TemplateField struct {
	Item *Item
	Type string
}

func GetTemplates(m ItemMap, fm FieldValueMap) []*Template {
	templateItem := m.FindItemByPath("/sitecore/templates/System/Templates/Template")

	list := m.FindItemsByTemplate(templateItem.ID)
	for _, item := range list {
		mapFields(item, fm)
	}
	return nil
}

func mapFields(item *ItemNode, fm FieldValueMap) {
	fmt.Println("Checking template", item.Name, item.ID, len(item.Children))
	for _, section := range item.Children {
		for _, flditem := range section.Children {
			fl, ok := fm[flditem.ID]
			if !ok {
				fmt.Println("no field values on ", flditem.ID)
				continue
			}

			for _, fv := range fl {
				if fv.Name == "Type" {
					fmt.Println("field", flditem.Name, " has type", fv.Value)
				}
			}
		}
	}
}
