package api

import (
	"fmt"
	"github.com/jasontconnell/sitecore/data"
)

func BuildUpdateItems(filteredMap data.ItemMap, referenceList []data.ItemNode, updateList []data.ItemNode) ([]data.UpdateItem, []data.UpdateField) {
	updateItems := []data.UpdateItem{}
	updateFields := []data.UpdateField{}
	itemMap := make(data.ItemMap)
	fieldMap := make(map[string]data.FieldValueNode)

	for _, sitem := range referenceList {
		itemMap[sitem.GetId()] = sitem
		for _, field := range sitem.GetFieldValues() {
			key := getFieldKey(sitem, field)
			fieldMap[key] = field
		}
	}

	deserializedItemMap := make(data.ItemMap)
	deserializedFieldMap := make(map[string]data.FieldValueNode)

	for _, ditem := range updateList {
		deserializedItemMap[ditem.GetId()] = ditem
		for _, dfield := range ditem.GetFieldValues() {
			key := getFieldKey(ditem, dfield)
			deserializedFieldMap[key] = dfield
		}
	}

	for _, sitem := range referenceList {
		_, inFilter := filteredMap[sitem.GetId()]
		if _, ok := deserializedItemMap[sitem.GetId()]; !ok && inFilter {
			updateItems = append(updateItems, data.UpdateItemFromItemNode(sitem, data.Delete))

			for _, sfield := range sitem.GetFieldValues() {
				updateFields = append(updateFields, data.UpdateFieldFromFieldValue(sfield, data.Delete))
			}
		}
	}

	for _, ditem := range updateList {
		if item, ok := itemMap[ditem.GetId()]; !ok {
			updateItems = append(updateItems, data.UpdateItemFromItemNode(ditem, data.Insert))
			for _, field := range ditem.GetFieldValues() {
				updateFields = append(updateFields, data.UpdateFieldFromFieldValue(field, data.Insert))
			}
		} else {
			for _, field := range ditem.GetFieldValues() {
				key := getFieldKey(ditem, field)
				if existingField, ok := fieldMap[key]; !ok {
					updateFields = append(updateFields, data.UpdateFieldFromFieldValue(existingField, data.Insert))
				} else {
					if existingField.GetValue() != field.GetValue() || existingField.GetVersion() != field.GetVersion() || existingField.GetLanguage() != field.GetLanguage() {
						updateFields = append(updateFields, data.UpdateFieldFromFieldValue(field, data.Update))
					}
				}
			}

			if item.GetName() != ditem.GetName() || item.GetTemplateId() != ditem.GetTemplateId() || item.GetMasterId() != ditem.GetMasterId() || item.GetParentId() != ditem.GetParentId() {
				updateItems = append(updateItems, data.UpdateItemFromItemNode(item, data.Update))
			}
		}
	}
	return updateItems, updateFields
}

func getFieldKey(item data.ItemNode, fv data.FieldValueNode) string {
	return fmt.Sprintf("%v_%v_%v", item.GetId(), fv.GetFieldId(), fv.GetLanguage())
}
