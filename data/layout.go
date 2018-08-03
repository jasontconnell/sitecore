package data

import (
	"github.com/google/uuid"
)

type RenderingType int

const (
	Sublayout RenderingType = iota
	Controller
	View
	XSLTBullshit
	SomeOtherCrap
)

type Rendering struct {
	Type RenderingType
	Item ItemNode
	Uid  uuid.UUID
}

type RenderingInstance struct {
	Rendering   Rendering
	Placeholder string
	Uid         uuid.UUID
	DataSource  string
}

type Device struct {
	Item ItemNode
}

type DeviceRendering struct {
	Device             Device
	RenderingInstances []RenderingInstance
}

type Layout struct {
	DeviceRenderings []DeviceRendering
}
