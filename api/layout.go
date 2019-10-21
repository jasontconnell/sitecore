package api

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/jasontconnell/sitecore/data"
	xr "github.com/jasontconnell/sitecore/data/xml"
)

func getRenderings(items []data.ItemNode, t data.RenderingType) []data.Rendering {
	all := []data.Rendering{}
	for _, item := range items {
		r := data.Rendering{Type: t, Item: item}
		switch t {
		case data.ControllerRenderingType:
			r.Info = getFieldValue("Controller", item.GetFieldValues()) + "." + getFieldValue("Controller Action", item.GetFieldValues())
		case data.SublayoutRenderingType, data.ViewRenderingType, data.LayoutRenderingType:
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

func parseRenderings(xmldata string) (xr.Root, error) {
	// check this occasionally. when merged (go 1.12?), the following if statement can probably be removed completely
	// https://go-review.googlesource.com/c/go/+/109855
	// there's a problem with the same object being used to parse essentially two versions of the same xml, one with and one without namespaces

	if strings.IndexAny(xmldata, ` s:id="{`) != -1 {
		xmldata = strings.Replace(xmldata, "s:ph=", "ph=", -1)
		xmldata = strings.Replace(xmldata, "s:id=", "id=", -1)
		xmldata = strings.Replace(xmldata, "s:ds=", "ds=", -1)
		xmldata = strings.Replace(xmldata, "s:par=", "par=", -1)
		xmldata = strings.Replace(xmldata, "p:after=", "after=", -1)
		xmldata = strings.Replace(xmldata, "p:before=", "before=", -1)
		xmldata = strings.Replace(xmldata, "<p:d />", "<d delete=\"true\" />", -1) // maybe this has to happen all the time, regardless of go bug mentioned above
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

func GetLayouts(m data.ItemMap) []data.Rendering {
	controllers := getRenderings(m.FindItemsByTemplate(data.ControllerRenderingId), data.ControllerRenderingType)
	sublayouts := getRenderings(m.FindItemsByTemplate(data.SublayoutRenderingId), data.SublayoutRenderingType)
	views := getRenderings(m.FindItemsByTemplate(data.ViewRenderingId), data.ViewRenderingType)
	webcontrols := getRenderings(m.FindItemsByTemplate(data.WebControlRenderingId), data.OtherRenderingType)
	xslcontrols := getRenderings(m.FindItemsByTemplate(data.XslRenderingId), data.OtherRenderingType)
	layouts := getRenderings(m.FindItemsByTemplate(data.LayoutTemplateId), data.LayoutRenderingType)

	return joinRenderings(controllers, sublayouts, views, webcontrols, xslcontrols, layouts)
}

type finalRenderingXml struct {
	xml     string
	version int64
}

func GetLayoutXml(item data.ItemNode) (renderings string, finalRenderings []finalRenderingXml) {
	if len(item.GetFieldValues()) == 0 {
		return renderings, finalRenderings
	}

	for _, fv := range item.GetFieldValues() {
		if fv.GetFieldId() == data.RenderingsFieldId {
			renderings = fv.GetValue()
		} else if fv.GetFieldId() == data.FinalRenderingsFieldId {
			finalRenderings = append(finalRenderings, finalRenderingXml{xml: fv.GetValue(), version: fv.GetVersion()})
		}
	}

	return renderings, finalRenderings
}

func MapAllLayouts(m data.ItemMap, tm data.TemplateMap, strict bool) error {
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
		var svRenderings string
		var svfinalRenderings []finalRenderingXml

		if t.HasStandardValues() {
			sv, svok := m[t.GetStandardValuesId()]
			if svok {
				svRenderings, svfinalRenderings = GetLayoutXml(sv)
			}
		}

		rr, err := getRenderingsFromXml(renderings, "Item", 1, m, rendmap, false, strict)
		if err != nil {
			return fmt.Errorf("Renderings from renderings, %v. %v", item.GetId(), err)
		}
		srr, err := getRenderingsFromXml(svRenderings, "Standard Values", 1, m, rendmap, true, strict)
		if err != nil {
			return fmt.Errorf("Renderings from standard value renderings, %v. %v", item.GetId(), err)
		}

		for _, r := range rr {
			item.AddRendering(r)
		}

		for _, sr := range srr {
			item.AddRendering(sr)
		}

		for _, frfld := range finalRenderings {
			fr, err := getRenderingsFromXml(frfld.xml, "Item", frfld.version, m, rendmap, false, strict)

			if err != nil {
				return fmt.Errorf("Renderings from final renderings, %v. %v", item.GetId(), err)
			}

			for _, r := range fr {
				item.AddFinalRendering(r)
			}
		}

		for _, frfld := range svfinalRenderings {
			sfr, err := getRenderingsFromXml(frfld.xml, "Standard Values", frfld.version, m, rendmap, true, strict)
			if err != nil {
				return fmt.Errorf("Renderings from standard values final renderings, %v. %v", item.GetId(), err)
			}

			for _, sr := range sfr {
				item.AddFinalRendering(sr)
			}
		}
	}

	return nil
}

func filterRenderingsOnItem(item data.ItemNode, rmap data.ItemMap, visited map[uuid.UUID]bool) {
	if _, ok := visited[item.GetId()]; ok {
		return
	}
	visited[item.GetId()] = true

	template := item.GetTemplate()

	idToUid := make(map[uuid.UUID]uuid.UUID) // map rendering id to uid

	removeUidMap := make(map[uuid.UUID]bool)

	tlist := append(template.GetBaseTemplates(), template)
	for _, t := range tlist {
		sv := t.GetStandardValues()
		if sv == nil {
			continue
		}
		filterRenderingsOnItem(sv, rmap, visited)

		for _, r := range append(sv.GetRenderings(), sv.GetFinalRenderings()...) {
			if r.StandardValues {
				continue
			}

			for _, rx := range r.RenderingInstances {
				if rx.Uid != emptyUuid && rx.Rendering.Item != nil {
					idToUid[rx.Rendering.Item.GetId()] = rx.Uid
				}
			}
		}
	}

	renderings := item.GetRenderings()
	toRemove := []data.RenderingInstance{}
	for _, drend := range renderings {
		for _, r := range drend.RenderingInstances {
			if r.Rendering.Item == nil { // some nodes won't have id so just move on.
				toRemove = append(toRemove, r)
				removeUidMap[r.Uid] = true
				continue
			}

			if _, ok := rmap[r.Rendering.Item.GetId()]; !ok {
				toRemove = append(toRemove, r)
				// renderings to remove by uid
				if r.Uid != emptyUuid {
					removeUidMap[r.Uid] = true
				}
				continue
			}

			if _, ok := removeUidMap[r.Uid]; ok && r.Uid != emptyUuid {
				toRemove = append(toRemove, r)
				continue
			}
		}
	}

	finalRenderings := item.GetFinalRenderings()

	for _, drend := range finalRenderings {
		for _, r := range drend.RenderingInstances {
			if r.Rendering.Item != nil {
				_, ok := rmap[r.Rendering.Item.GetId()]
				if !ok {
					toRemove = append(toRemove, r)
					if r.Uid != emptyUuid {
						removeUidMap[r.Uid] = true
					}
				}
			}

			if _, ok := removeUidMap[r.Uid]; ok && r.Uid != emptyUuid {
				toRemove = append(toRemove, r)
			}
		}
	}

	for _, rem := range toRemove {
		item.RemoveRendering(rem)
		item.RemoveFinalRendering(rem)
	}
}

// rmap is the included renderings
func FilterRenderings(itemMap data.ItemMap, rmap data.ItemMap) {
	visited := make(map[uuid.UUID]bool)
	for _, m := range itemMap {
		filterRenderingsOnItem(m, rmap, visited)
	}
}

func FilterDataSources(itemMap data.ItemMap) {
	for _, item := range itemMap {
		for _, dr := range item.GetRenderings() {
			for _, r := range dr.RenderingInstances {
				if _, ok := itemMap[r.DataSourceId]; !ok {
					r.DataSourceId = data.EmptyID
					r.DataSource = ""
				}
			}
		}

		for _, dr := range item.GetFinalRenderings() {
			for _, r := range dr.RenderingInstances {
				if _, ok := itemMap[r.DataSourceId]; !ok {
					r.DataSourceId = data.EmptyID
					r.DataSource = ""
				}
			}
		}
	}
}

func UpdateRenderingsFields(itemMap data.ItemMap) {
	for _, item := range itemMap {
		UpdateRenderingsXml(item)
		UpdateFinalRenderingsXml(item)
	}
}

func UpdateRenderingsXml(item data.ItemNode) {
	x := xr.SerializeRenderings(item.GetRenderings())
	fields := item.GetFieldValues()
	for _, f := range fields {
		if f.GetFieldId() == data.RenderingsFieldId {
			f.SetValue(x)
			break
		}
	}
}

func UpdateFinalRenderingsXml(item data.ItemNode) {
	// get renderings xml
	// add namespaces to selected attributes (before, after, id, par, etc)
	x := xr.SerializeFinalRenderings(item.GetFinalRenderings())
	fields := item.GetFieldValues()
	for _, f := range fields {
		if f.GetFieldId() == data.FinalRenderingsFieldId {
			f.SetValue(x)
			break
		}
	}
}

func getRenderingsFromXml(x, loc string, version int64, m data.ItemMap, rendmap map[uuid.UUID]data.Rendering, standardValues, strict bool) ([]data.DeviceRendering, error) {
	if len(x) == 0 {
		return []data.DeviceRendering{}, nil
	}

	rs, rerr := parseRenderings(x)
	if rerr != nil {
		return nil, rerr
	}

	drends := []data.DeviceRendering{}
	for _, dx := range rs.Devices {
		devid := MustParseUUID(dx.ID)
		deviceItem, ok := m[devid]
		if !ok && strict {
			return nil, fmt.Errorf("Can't find item with id for device, %v", dx.ID)
		}

		var layoutid = data.EmptyID
		var layout data.Layout
		if len(dx.Layout) != 0 {
			layoutid = MustParseUUID(dx.Layout)
			layoutItem, lok := m[layoutid]

			if lok {
				layout.Path = getFieldValue("Path", layoutItem.GetFieldValues())
				layout.Item = layoutItem
			}
		}

		device := data.Device{Item: deviceItem, Layout: layout}

		dr := data.DeviceRendering{Device: device, Version: version, RenderingInstances: []data.RenderingInstance{}, StandardValues: standardValues}

		for _, rx := range dx.Renderings {
			// rendering id can be not provided.
			rid, err := TryParseUUID(rx.ID)
			ruid := MustParseUUID(rx.Uid)

			rend, ok := rendmap[rid]
			if !ok && !rx.DeleteNode.Delete && rid != emptyUuid && strict {
				return nil, fmt.Errorf("Couldn't find rendering with id %v in %v", rid, x)
			}

			dsid, err := TryParseUUID(rx.DataSource)

			if err != nil { // not a uuid or blank
				rx.DataSource = ""
			}

			rinst := data.RenderingInstance{Rendering: rend, Placeholder: rx.Placeholder, Uid: ruid, DataSource: rx.DataSource, DataSourceId: dsid, Deleted: rx.DeleteNode.Delete}
			if len(rx.Parameters) > 0 {
				parstr, err := url.PathUnescape(rx.Parameters)
				if err != nil {
					return nil, fmt.Errorf("Couldn't path unescape %s.  %v", rx.Parameters, err)
				}
				params := strings.Split(parstr, "&") // after xml entity ref processing it's just & not &amp;
				for _, p := range params {
					ps := strings.Split(p, "=")
					k := ps[0]
					v := ""
					if len(ps) > 1 {
						v = ps[1]
					}
					rinst.Parameters = append(rinst.Parameters, data.KV{Key: k, Value: v})
				}
			}

			dr.RenderingInstances = append(dr.RenderingInstances, rinst)
		}

		drends = append(drends, dr)
	}

	return drends, nil
}

func GetRenderingItems(m data.ItemMap) data.ItemMap {
	tmap := make(map[uuid.UUID]bool)
	renderingItems := data.ItemMap{}

	for _, id := range []uuid.UUID{
		data.ControllerRenderingId,
		data.ItemRenderingId,
		data.MethodRenderingId,
		data.SublayoutRenderingId,
		data.UrlRenderingId,
		data.ViewRenderingId,
		data.WebControlRenderingId,
		data.XmlControlRenderingId,
		data.XslRenderingId,
		data.LayoutTemplateId,
	} {
		tmap[id] = true
	}

	for _, item := range m {
		if _, ok := tmap[item.GetTemplateId()]; ok {
			renderingItems[item.GetId()] = item
		}
	}

	return renderingItems
}
