package data

import (
	"github.com/google/uuid"
)

type UpdateType string

const (
	Insert         UpdateType = "insert"
	Update         UpdateType = "update"
	Delete         UpdateType = "delete"
	Ignore         UpdateType = "ignore"
	InsertOrUpdate UpdateType = "insertorupdate"
)

type UpdateItem struct {
	ID         uuid.UUID
	Name       string
	TemplateID uuid.UUID
	ParentID   uuid.UUID
	MasterID   uuid.UUID
	UpdateType UpdateType
}

type UpdateField struct {
	ItemID     uuid.UUID
	FieldID    uuid.UUID
	Value      string
	Source     FieldSource
	Version    int64
	Language   string
	UpdateType UpdateType
}

func UpdateItemFromItemNode(node ItemNode, updateType UpdateType) UpdateItem {
	item := UpdateItem{}
	item.ID = node.GetId()
	item.Name = node.GetName()
	item.TemplateID = node.GetTemplateId()
	item.ParentID = node.GetParentId()
	item.MasterID = node.GetMasterId()
	item.UpdateType = updateType

	return item
}

func UpdateFieldFromFieldValue(fv FieldValueNode, updateType UpdateType) UpdateField {
	fld := UpdateField{}
	fld.ItemID = fv.GetItemId()
	fld.FieldID = fv.GetFieldId()
	fld.Value = fv.GetValue()
	fld.Source = fv.GetSource()
	fld.Version = fv.GetVersion()
	fld.Language = fv.GetLanguage()
	fld.UpdateType = updateType
	return fld
}
