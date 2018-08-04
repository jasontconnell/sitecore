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

func getRenderings(items []data.ItemNode, t data.RenderingType) []data.Rendering {
	all := []data.Rendering{}
	for _, item := range items {
		r := data.Rendering{Type: t, Item: item}
		switch t {
		case data.Controller:
			r.Info = getFieldValue("Controller", item.GetFieldValues()) + "." + getFieldValue("Controller Action", item.GetFieldValues())
		case data.Sublayout,data.View:
			r.Info = getFieldValue("Path", item.GetFieldValues())
		}
		all = append(all, r)
	}
	return all
}

func joinRenderings(lists ...[]data.Rendering) []data.Rendering {
	all := []data.Rendering{}

	for _, list := range lists {
		for _, r := range list {
			all = append(all, r)
		}

	}
	return all
}

func GetLayouts(m data.ItemMap) []data.Rendering {
	controllers := getRenderings(m.FindItemsByTemplate(data.ControllerRenderingId), data.Controller)
	sublayouts := getRenderings(m.FindItemsByTemplate(data.SublayoutRenderingId), data.Sublayout)
	views := getRenderings(m.FindItemsByTemplate(data.ViewRenderingId), data.View)
	webcontrols := getRenderings(m.FindItemsByTemplate(data.WebControlRenderingId), data.Other)
	xslcontrols := getRenderings(m.FindItemsByTemplate(data.XslRenderingId), data.Other)

	return joinRenderings(controllers, sublayouts, views, webcontrols, xslcontrols)
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

		var layoutid = data.EmptyID
		var layout data.ItemNode
		if len(dx.Layout) != 0 {
			layoutid = MustParseUUID(dx.Layout)
			layout, _ = m[layoutid]
		}

		device := data.Device{Item: deviceItem, Layout: layout}

		dr := data.DeviceRendering{Device: device, RenderingInstances: []data.RenderingInstance{}}

		for _, rx := range dx.Renderings {
			if rx.ID == "" {
				continue
			}
			rid := MustParseUUID(rx.ID)
			ruid := MustParseUUID(rx.Uid)

			rend, ok := rendmap[rid]
			if !ok {
				rend = data.Rendering{ ID: rid, Type: data.NotFound }
				rendmap[rid] = rend
				// return nil, fmt.Errorf("Couldn't find rendering with id %v in %v", rid, x)
				// fmt.Printf("Couldn't find rendering with id %v in %v\n", rid, x)
				// fmt.Printf("Item is %v\n", m[rid])
				// continue
			}

			rinst := data.RenderingInstance{Rendering: rend, Placeholder: rx.Placeholder, Uid: ruid, DataSource: rx.DataSource}
			dr.RenderingInstances = append(dr.RenderingInstances, rinst)
		}

		drends = append(drends, dr)
	}

	return drends, nil
}
