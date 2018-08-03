package api

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"github.com/jasontconnell/sitecore/data"
	xr "github.com/jasontconnell/sitecore/data/xml"
	"io"
	"strings"
)

func joinRenderings(l ...[]data.ItemNode) []data.Rendering {
	all := []data.Rendering{}

	for _, list := range l {
		for _, item := range list {
			r := data.Rendering{Type: data.Controller, Item: item}
			all = append(all, r)
		}

	}
	return all
}

func GetLayouts(m data.ItemMap) []data.Rendering {
	controllers := m.FindItemsByTemplate(data.ControllerRenderingId)
	sublayouts := m.FindItemsByTemplate(data.SublayoutRenderingId)
	views := m.FindItemsByTemplate(data.ViewRenderingId)
	webcontrols := m.FindItemsByTemplate(data.WebControlRenderingId)

	return joinRenderings(controllers, sublayouts, views, webcontrols)
}

func GetLayoutXml(item data.ItemNode) (renderings, finalRenderings string) {
	if len(item.GetFieldValues()) == 0 {
		return renderings, finalRenderings
	}

	for _, fv := range item.GetFieldValues() {
		if fv.GetFieldId() == data.RenderingsFieldId {
			renderings = fv.GetValue()
		} else if fv.GetFieldId() == data.FinalRenderingsFieldId {
			finalRenderings = fv.GetValue()
		}
	}

	return renderings, finalRenderings
}

func GetRenderings(xmldata string) (xr.Root, error) {
	if strings.IndexAny(xmldata, ` s:id="{`) != -1 {
		xmldata = strings.Replace(strings.Replace(xmldata, "s:ph=", "ph=", -1), "s:id=", "id=", -1)
	}

	b := bytes.NewBufferString(xmldata)

	r := xr.Root{}

	dec := xml.NewDecoder(b)
	dec.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		return input, nil
	}
	xerr := dec.Decode(&r)

	return r, xerr
}

func MapAllLayouts(m data.ItemMap, tm map[uuid.UUID]data.TemplateNode) error {
	rendlist := GetLayouts(m)
	rendmap := make(map[uuid.UUID]data.Rendering)

	for _, r := range rendlist {
		rendmap[r.Item.GetId()] = r
	}

	for _, item := range m {
		fmt.Println("doing it for ", item.GetId(), item.GetPath())
		t, ok := tm[item.GetTemplateId()]
		if !ok {
			return fmt.Errorf("Template not found, %v", item.GetTemplateId())
		}

		renderings, finalRenderings := GetLayoutXml(item)
		var svrenderings, svfinalRenderings string

		if t.HasStandardValues() {
			sv, svok := m[t.GetStandardValuesId()]
			if svok {
				svrenderings, svfinalRenderings = GetLayoutXml(sv)
			}
		}

		rr, err := getRenderingsFromXml(renderings, "Item", m, rendmap)
		if err != nil {
			return fmt.Errorf("Renderings from renderings, %v", err)
		}
		srr, err := getRenderingsFromXml(svrenderings, "Standard Values", m, rendmap)
		if err != nil {
			return fmt.Errorf("Renderings from standard value renderings, %v", err)
		}

		for _, r := range rr {
			item.AddRendering(r)
		}

		for _, sr := range srr {
			item.AddRendering(sr)
		}

		fr, err := getRenderingsFromXml(finalRenderings, "Item", m, rendmap)
		if err != nil {
			return fmt.Errorf("Renderings from final renderings, %v", err)
		}

		sfr, err := getRenderingsFromXml(svfinalRenderings, "Standard Values", m, rendmap)
		if err != nil {
			return fmt.Errorf("Renderings from standard values final renderings, %v", err)
		}
		for _, r := range fr {
			item.AddFinalRendering(r)
		}

		for _, sr := range sfr {
			item.AddFinalRendering(sr)
		}
	}

	return nil
}

func getRenderingsFromXml(x, loc string, m data.ItemMap, rendmap map[uuid.UUID]data.Rendering) ([]data.DeviceRendering, error) {
	if len(x) == 0 {
		return []data.DeviceRendering{}, nil
	}

	rs, rerr := GetRenderings(x)
	if rerr != nil {
		return nil, rerr
	}

	drends := []data.DeviceRendering{}
	for _, dx := range rs.Devices {
		devid := MustParseUUID(dx.ID)
		deviceItem, ok := m[devid]
		if !ok {
			return nil, fmt.Errorf("Can't find item with id for device, %v", dx.ID)
		}
		device := data.Device{Item: deviceItem}

		dr := data.DeviceRendering{Device: device, RenderingInstances: []data.RenderingInstance{}}

		for _, rx := range dx.Renderings {
			if rx.ID == "" {
				continue
			}
			rid := MustParseUUID(rx.ID)
			ruid := MustParseUUID(rx.Uid)

			rend, ok := rendmap[rid]
			if !ok {
				//return nil, fmt.Errorf("Couldn't find rendering with id %v in %v", rid, x)
				fmt.Printf("Couldn't find rendering with id %v in %v\n", rid, x)
				continue
			}

			rinst := data.RenderingInstance{Rendering: rend, Placeholder: rx.Placeholder, Uid: ruid, DataSource: rx.DataSource}
			dr.RenderingInstances = append(dr.RenderingInstances, rinst)
		}

		drends = append(drends, dr)
	}

	return drends, nil
}
