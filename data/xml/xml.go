package xml

import (
	"encoding/xml"
)

type Root struct {
	XMLName xml.Name `xml:"r"`
	Devices []Device `xml:"d"`
}

// type FinalRenderingsRoot struct {
// 	XMLName xml.Name `xml:"r"`
// 	Devices []FinalRenderingsDevice `xml:"d"`
// }

// type FinalRenderingsDevice struct {
// 	XMLName      xml.Name      `xml:"d"`
// 	ID           string        `xml:"id,attr"`
// 	Placeholders []FinalRenderingsPlaceholder `xml:"p"`
// 	Renderings   []FinalRendering   `xml:"r"`
// }

// type FinalRenderingsPlaceholder struct {
// 	XMLName xml.Name `xml:"p"`
// 	ID      string   `xml:"uid,attr"`
// }

// type FinalRendering struct {
// 	XMLName     xml.Name `xml:"r"`
// 	ID string `xml:"s id,attr"`
// 	Uid         string   `xml:"uid,attr"`
// 	Before      string   `xml:"p before,attr"`
// 	Placeholder string   `xml:"s ph,attr"`
// 	DataSource string `xml:"s ds,attr"`
// }

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
	XMLName xml.Name `xml:"r"`
	ID      string   `xml:"id,attr"`
	//SID string `xml:"s id,attr"`
	Uid         string `xml:"uid,attr"`
	Placeholder string `xml:"ph,attr"`
	//SPlaceholder string `xml:"s ph,attr"`
	DataSource string `xml:"ds,attr"`
}
