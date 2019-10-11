package data

import (
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

type item struct {
	ID         uuid.UUID
	Name       string
	TemplateID uuid.UUID
	ParentID   uuid.UUID
	MasterID   uuid.UUID
	Created    time.Time
	Updated    time.Time
	Path       string
	Parent     ItemNode
	Children   []ItemNode
	Level      int

	FieldValues []FieldValueNode

	Renderings      []DeviceRendering
	FinalRenderings []DeviceRendering

	Template TemplateNode

	versions map[int64]int64
}

func NewBlankItemNode() ItemNode {
	return new(item)
}

func NewItemNode(id uuid.UUID, name string, templateId, parentId, masterId uuid.UUID) ItemNode {
	item := &item{}
	item.ID = id
	item.Name = name
	item.TemplateID = templateId
	item.ParentID = parentId
	item.MasterID = masterId
	return item
}

func (item *item) GetId() uuid.UUID {
	return item.ID
}

func (item *item) SetId(id uuid.UUID) {
	item.ID = id
}

func (item *item) GetChildren() []ItemNode {
	return item.Children
}

func (item *item) AddChild(node ItemNode) {
	item.Children = append(item.Children, node)
}

func (item *item) GetLevel() int {
	return item.Level
}

func (item *item) SetLevel(level int) {
	item.Level = level
}

func (item *item) GetPath() string {
	return item.Path
}

func (item *item) SetPath(p string) {
	item.Path = p
}

func (item *item) GetName() string {
	return item.Name
}

func (item *item) SetName(n string) {
	item.Name = n
}

func (item *item) GetParentId() uuid.UUID {
	return item.ParentID
}

func (item *item) SetParentId(id uuid.UUID) {
	item.ParentID = id
}

func (item *item) GetParent() ItemNode {
	return item.Parent
}

func (item *item) SetParent(node ItemNode) {
	item.Parent = node
}

func (item *item) GetMasterId() uuid.UUID {
	return item.MasterID
}

func (item *item) SetMasterId(id uuid.UUID) {
	item.MasterID = id
}

func (item *item) GetTemplateId() uuid.UUID {
	return item.TemplateID
}

func (item *item) SetTemplateId(id uuid.UUID) {
	item.TemplateID = id
}

func (item *item) GetTemplate() TemplateNode {
	return item.Template
}

func (item *item) SetTemplate(t TemplateNode) {
	item.Template = t
}

func (t *item) String() string {
	return fmt.Sprintf("ID: %v\nName:%v\nPath:%v\n", t.ID, t.Name, t.Path)
}

func (t *item) GetFieldValues() []FieldValueNode {
	return t.FieldValues
}

func (t *item) AddFieldValue(fv FieldValueNode) {
	if _, ok := t.versions[fv.GetVersion()]; !ok && fv.GetSource() == VersionedFields {
		t.versions[fv.GetVersion()] = fv.GetVersion()
	}

	t.FieldValues = append(t.FieldValues, fv)
}

func (t item) GetVersions() []int64 {
	list := []int64{}
	for _, v := range t.versions {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})
	return list
}

func (t *item) GetRenderings() []DeviceRendering {
	return t.Renderings
}

func (t *item) GetFinalRenderings() []DeviceRendering {
	v := t.GetVersions()
	var vnum int64 = 1
	if len(v) > 0 {
		vnum = v[len(v)-1]
	}

	frens := []DeviceRendering{}
	for _, dr := range t.FinalRenderings {
		if dr.Version == vnum {
			frens = append(frens, dr)
		}
	}

	return frens
}

func (t *item) AddRendering(r DeviceRendering) {
	t.Renderings = append(t.Renderings, r)
}

func (t *item) AddFinalRendering(r DeviceRendering) {
	t.FinalRenderings = append(t.FinalRenderings, r)
}

func (t *item) RemoveRendering(r RenderingInstance) {
	res := []DeviceRendering{}
	for _, ritr := range t.Renderings {
		rs := []RenderingInstance{}
		for _, iri := range ritr.RenderingInstances {
			if iri.Uid != r.Uid {
				rs = append(rs, iri)
			}
		}
		ritr.RenderingInstances = rs

		res = append(res, ritr)
	}

	t.Renderings = res
}

func (t *item) RemoveFinalRendering(r RenderingInstance) {
	res := []DeviceRendering{}
	for _, ritr := range t.FinalRenderings {
		rs := []RenderingInstance{}
		for _, iri := range ritr.RenderingInstances {
			if iri.Uid != r.Uid {
				rs = append(rs, iri)
			}
		}
		ritr.RenderingInstances = rs
		res = append(res, ritr)
	}

	t.FinalRenderings = res
}
