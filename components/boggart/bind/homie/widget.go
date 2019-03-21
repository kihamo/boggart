package homie

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	var err error

	if r.IsPost() {
		switch r.URL().Query().Get("action") {
		case "reset":
			err = bind.Reset(r.Context())
			if err != nil {
				w.SendJSON(response{
					Result:  "failed",
					Message: err.Error(),
				})

			} else {
				w.SendJSON(response{
					Result:  "success",
					Message: "Send reset signal success",
				})
			}

		case "config":
			err = r.Original().ParseForm()
			if err == nil {
				for key, value := range r.Original().PostForm {
					if len(value) == 0 {
						continue
					}

					err = bind.ImplementationConfigSet(r.Context(), key, value[0])
					if err != nil {
						break
					}
				}
			}

			if err != nil {
				w.SendJSON(response{
					Result:  "failed",
					Message: err.Error(),
				})

			} else {
				w.SendJSON(response{
					Result:  "success",
					Message: "Send signal of set config success",
				})
			}

		default:
			t.NotFound(w, r)
		}

		return
	}

	vars := map[string]interface{}{
		"error":              err,
		"name":               "",
		"devices_attributes": bind.DeviceAttributes(),
		"config":             bind.ImplementationConfigAll(),
	}

	if attribute, ok := bind.DeviceAttribute("name"); ok {
		vars["name"] = attribute
	}

	t.Render(r.Context(), "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
