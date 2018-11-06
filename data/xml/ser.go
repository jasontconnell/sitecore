package xml

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jasontconnell/sitecore/data"
	"strings"
)

const (
	root = `<r xmlns:xsd="http://www.w3.org/2001/XMLSchema">
	%s
</r>`
	device = `	<d id="%s" l="%s">
		%s
	</d>`

	placeholder = `		<p uid="%s" key="%s" md="%s" />`

	rootFinal = `<r xmlns:p="p" xmlns:s="s" p:p="1">
	%s
</r>`

	deviceFinal = `	<d id="%s">
		%s
	</d>`

	renderingFinal        = `		<r uid="%s" %s%ss:par="%s" s:ph="%s" />`
	renderingFinalDeleted = `		<r uid="%s">
			<p:d />
		</r>`
)

func fmtxml(x string) string {
	return strings.Replace(strings.Replace(x, "\t", "", -1), "\n", "", -1)
}

func SerializeRenderings(ren []data.DeviceRendering) string {
	dx := ""
	for _, dv := range ren {
		if dv.StandardValues {
			continue
		}
		phx := ""
		for _, p := range dv.PlaceholderInstances {
			phx += fmt.Sprintf(placeholder, p.Uid, p.Key, p.PlaceholderPath)
		}

		rx := ""
		for _, r := range dv.RenderingInstances {
			rx += `<r `
			if r.DataSource != "" {
				rx += fmt.Sprintf(`ds="%s" `, r.DataSource)
			}

			if r.Rendering.Item != nil {
				rx += fmt.Sprintf(`id="%s" `, idFormat(r.Rendering.Item.GetId()))
			}

			par := getParamString(r.Parameters)
			if par != "" {
				rx += fmt.Sprintf(`par="%s" `, par)
			}

			if r.Uid != data.EmptyID {
				rx += fmt.Sprintf(`uid="%s" `, idFormat(r.Uid))
			}

			if r.Placeholder != "" {
				rx += fmt.Sprintf(`ph="%s" `, r.Placeholder)
			}

			rx += " />"
			//rx += fmt.Sprintf(rendering, r.DataSource, idFormat(r.Rendering.Item.GetId()), getParamString(r.Parameters), idFormat(r.Uid), r.Placeholder)
		}

		var layoutid string
		if dv.Device.Layout.Item != nil {
			layoutid = idFormat(dv.Device.Layout.Item.GetId())
		}
		dx += fmt.Sprintf(device, idFormat(dv.Device.Item.GetId()), layoutid, phx+rx)
	}

	// remove whitespace and new lines
	return fmtxml(fmt.Sprintf(root, dx))
}

func SerializeFinalRenderings(ren []data.DeviceRendering) string {
	dx := ""
	for _, dv := range ren {
		if dv.StandardValues {
			continue
		}
		phx := ""
		for _, p := range dv.PlaceholderInstances {
			phx += fmt.Sprintf(placeholder, p.Uid, p.Key, p.PlaceholderPath)
		}

		rx := ""
		for _, r := range dv.RenderingInstances {
			if r.Deleted {
				rx += fmt.Sprintf(renderingFinalDeleted, idFormat(r.Uid))
			} else {

				rx += "<r "
				if r.DataSource != "" {
					rx += fmt.Sprintf(`s:ds="%s" `, r.DataSource)
				}

				if r.Before != "" || r.After != "" {
					val := r.Before
					befaft := "before"
					if r.After != "" {
						befaft = "after"
						val = r.After
					}

					rx += fmt.Sprintf(`p:%s="%s" `, befaft, val)
				}

				if r.Rendering.Item != nil {
					rx += fmt.Sprintf(`s:id="%s" `, idFormat(r.Rendering.Item.GetId()))
				}

				par := getParamString(r.Parameters)
				if par != "" {
					rx += fmt.Sprintf(`s:par="%s" `, par)
				}

				if r.Uid != data.EmptyID {
					rx += fmt.Sprintf(`uid="%s" `, idFormat(r.Uid))
				}

				if r.Placeholder != "" {
					rx += fmt.Sprintf(`s:ph="%s" `, r.Placeholder)
				}

				rx += " />"

				//rx += fmt.Sprintf(renderingFinal, idFormat(r.Uid), befaftxml, idstr, r.DataSource, getParamString(r.Parameters), r.Placeholder)
			}

		}

		dx += fmt.Sprintf(deviceFinal, idFormat(dv.Device.Item.GetId()), phx+rx)
	}

	return fmtxml(fmt.Sprintf(rootFinal, dx))
}

func idFormat(uid uuid.UUID) string {
	s := uid.String()
	s = strings.ToUpper(s)
	return "{" + s + "}"
}

func getParamString(params []data.KV) string {
	s := ""
	for _, p := range params {
		s += p.Key + "=" + p.Value + "&"
	}
	return strings.TrimSuffix(s, "&")
}
