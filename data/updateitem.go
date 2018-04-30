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
	Source     string
	Version    int64
	Language   string
	UpdateType UpdateType
}

func UpdateItemFromItemNode(node ItemNode) UpdateItem {
	item := UpdateItem{}
	item.ID = node.GetId()
	item.Name = node.GetName()
	item.TemplateID = node.GetTemplateId()
	item.ParentID = node.GetParentId()
	item.MasterID = node.GetMasterId()
	item.UpdateType = Update

	return item
}

func UpdateFieldFromFieldValue(fv FieldValue) UpdateField {
	fld := UpdateField{}
	fld.ItemID = fv.ItemID
	fld.FieldID = fv.FieldID
	fld.Value = fv.Value
	fld.Source = fv.Source
	fld.Version = fv.Version
	fld.Language = fv.Language
	return fld
}
