package xml

import (
	"encoding/xml"
)

type Root struct {
	XMLName xml.Name `xml:"r"`
	Devices []Device `xml:"d"`
}

type Device struct {
	XMLName      xml.Name      `xml:"d"`
	ID           string        `xml:"id,attr"`
	Layout       string        `xml:"l,attr"`
	Placeholders []Placeholder `xml:"p"`
	Renderings   []Rendering   `xml:"r"`
}

type Placeholder struct {
	XMLName xml.Name `xml:"p"`
	ID      string   `xml:"uid,attr"`
}

type Rendering struct {
	XMLName     xml.Name `xml:"r"`
	ID          string   `xml:"id,attr"`
	Uid         string   `xml:"uid,attr"`
	Placeholder string   `xml:"ph,attr"`
	DataSource  string   `xml:"ds,attr"`
	Parameters  string   `xml:"par,attr"`
}
