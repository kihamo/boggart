package openhab

import (
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/openhab/client/items"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)

	action := r.URL().Query().Get("action")
	switch action {
	case "input":
		id := strings.TrimSpace(r.URL().Query().Get("id"))
		if id == "" {
			t.NotFound(w, r)
			return
		}

		if r.IsPost() {
			value := r.Original().FormValue("value")
			if value == "" {
				value = "NULL"
			}

			paramsPut := items.NewPutItemStateParams().
				WithItemname(id).
				WithBody(value)

			_, err := bind.provider.Items.PutItemState(paramsPut)
			if err != nil {
				t.InternalError(w, r, err)
				return
			}

			return
		}

		paramsGet := items.NewGetItemDataParams().
			WithItemname(id)

		response, err := bind.provider.Items.GetItemData(paramsGet)
		if err != nil {
			t.NotFound(w, r)
			return
		}

		style := r.URL().Query().Get("style")
		switch style {
		case "android", "basicui":
			// skip
		default:
			if strings.Contains(strings.ToLower(r.Original().Header.Get("User-Agent")), "android") {
				style = "android"
			} else {
				style = "basicui"
			}
		}

		if response.Payload.Type == "DateTime" {
			t, _ := time.Parse("2006-01-02T15:04:05.999Z0700", response.Payload.State)
			response.Payload.State = t.Format("2006-01-02T15:04:05")
		}

		var iconURL string

		if response.Payload.Category != "" {
			iconURL = r.URL().Path + "/?action=icon&icon=" + response.Payload.Category +
				"&state=" + response.Payload.State + "&format=svg&anyFormat=true"

			if key := r.URL().Query().Get(boggart.AccessKeyName); key != "" {
				iconURL += "&" + boggart.AccessKeyName + "=" + key
			}
		}

		t.RenderLayout(r.Context(), "input", "input", map[string]interface{}{
			"item":     response.Payload,
			"theme":    r.URL().Query().Get("theme"),
			"type":     r.URL().Query().Get("type"),
			"rows":     r.URL().Query().Get("rows"),
			"style":    style,
			"icon_url": iconURL,
		})
		return

	case "icon":
		req := r.Original().Clone(r.Context())
		req.URL.Path = "icon/" + r.URL().Query().Get("icon")

		bind.proxy.ServeHTTP(w, req)
		return
	}

	t.NotFound(w, r)
	return
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
