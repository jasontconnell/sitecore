package data

import (
	"time"

	"github.com/google/uuid"
)

type FieldValueKey struct {
	VersionKey
	FieldId uuid.UUID
}

type VersionKey struct {
	Language Language
	Version  int64
	Source   FieldSource
}

type ItemData interface {
	GetId() uuid.UUID
	SetId(id uuid.UUID)

	GetParentId() uuid.UUID
	SetParentId(id uuid.UUID)

	GetTemplateId() uuid.UUID
	SetTemplateId(id uuid.UUID)

	SetTemplate(t TemplateNode)
	GetTemplate() TemplateNode

	GetMasterId() uuid.UUID
	SetMasterId(id uuid.UUID)

	GetCreated() time.Time
	GetUpdated() time.Time

	GetLevel() int
	SetLevel(level int)
	GetPath() string
	SetPath(p string)

	GetName() string
	SetName(n string)

	GetRenderings() []DeviceRendering
	GetFinalRenderings() []DeviceRendering

	AddRendering(r DeviceRendering)
	AddFinalRendering(r DeviceRendering)

	RemoveRendering(r RenderingInstance)
	RemoveFinalRendering(r RenderingInstance)
}

type ItemNode interface {
	ItemData

	GetChildren() []ItemNode
	AddChild(node ItemNode)

	GetParent() ItemNode
	SetParent(node ItemNode)

	GetFieldValues() []FieldValueNode
	AddFieldValue(fv FieldValueNode)
	GetFieldValue(fieldId uuid.UUID, language Language) FieldValueNode

	GetVersionedFieldKeys(language Language, version int64) []FieldValueKey
	GetLatestVersionFieldKeys(language Language) []FieldValueKey
	GetVersionedFieldValues() map[FieldValueKey]FieldValueNode
	GetLatestVersionFields(language Language) []FieldValueNode
	GetFieldsByVersion(language Language, version int64) []FieldValueNode

	GetVersions() []int64
}

type ItemMap map[uuid.UUID]ItemNode

func (m ItemMap) FindItems(name string) []ItemNode {
	items := []ItemNode{}
	for _, item := range m {
		if item.GetName() == name {
			items = append(items, item)
		}
	}
	return items
}

func (m ItemMap) FindItemByPath(path string) ItemNode {
	var node ItemNode
	for _, item := range m {
		if item.GetPath() == path {
			node = item
		}
	}
	return node
}

func (m ItemMap) FindItemsByTemplate(templateId uuid.UUID) []ItemNode {
	return m.FindItemsByTemplates([]uuid.UUID{templateId})
}

func (m ItemMap) FindItemsByTemplates(templateIds []uuid.UUID) []ItemNode {
	list := []ItemNode{}
	tmap := make(map[uuid.UUID]bool)
	for _, tid := range templateIds {
		tmap[tid] = true
	}

	for _, item := range m {
		if _, ok := tmap[item.GetTemplateId()]; ok {
			list = append(list, item)
		}
	}

	return list
}
