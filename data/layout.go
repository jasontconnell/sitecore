package data

import (
	"github.com/google/uuid"
)

type RenderingType int

const (
	Sublayout RenderingType = iota
	Controller
	View
	Other
	NotFound
)

type Rendering struct {
	ID   uuid.UUID
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
	Rendering   Rendering
	Placeholder string
	Uid         uuid.UUID
	DataSource  string
	Parameters  []KV
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
	Device             Device
	RenderingInstances []RenderingInstance
}
