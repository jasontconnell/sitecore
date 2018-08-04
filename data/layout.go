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

type RenderingInstance struct {
	Rendering   Rendering
	Placeholder string
	Uid         uuid.UUID
	DataSource  string
}

type Device struct {
	Item   ItemNode
	Layout ItemNode
}

type DeviceRendering struct {
	Device             Device
	RenderingInstances []RenderingInstance
}
