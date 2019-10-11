package data

import (
	"fmt"

	"github.com/google/uuid"
)

type RenderingType int

const (
	SublayoutRenderingType RenderingType = iota
	ControllerRenderingType
	ViewRenderingType
	LayoutRenderingType
	OtherRenderingType
	NotFoundRenderingType
)

type Rendering struct {
	Type RenderingType
	Item ItemNode
	Uid  uuid.UUID
	Info string
}

type KV struct {
	Key   string
	Value string
}

type RenderingInstance struct {
	Rendering    Rendering
	Placeholder  string
	Uid          uuid.UUID
	DataSource   string
	DataSourceId uuid.UUID
	Parameters   []KV
	Before       string
	After        string
	Deleted      bool
	Modified     bool
}

type PlaceholderInstance struct {
	Uid             string
	Key             string
	PlaceholderPath string
}

type Device struct {
	Item   ItemNode
	Layout Layout
}

type Layout struct {
	Path string
	Item ItemNode
}

type DeviceRendering struct {
	Device               Device
	RenderingInstances   []RenderingInstance
	PlaceholderInstances []PlaceholderInstance
	StandardValues       bool
	Version              int64
}

func (dr DeviceRendering) String() string {
	return fmt.Sprintf("device: %s  renderings: %v  placeholder %v  version %v", dr.Device, dr.RenderingInstances, dr.PlaceholderInstances, dr.Version)
}

func (ri RenderingInstance) String() string {
	return fmt.Sprintf("rendering: %s ph: %s", ri.Rendering, ri.Placeholder)
}

func (r Rendering) String() string {
	var id uuid.UUID
	if r.Item != nil {
		id = r.Item.GetId()
	}
	return fmt.Sprintf("id: %s", id)
}

func (d Device) String() string {
	return fmt.Sprintf("device id: %s layout: %s", d.Item.GetId(), d.Layout)
}

func (ph PlaceholderInstance) String() string {
	return fmt.Sprintf("uid: %s key: %s path: %s", ph.Uid, ph.Key, ph.PlaceholderPath)
}

func (l Layout) String() string {
	return fmt.Sprintf("layout path: %s", l.Path)
}
