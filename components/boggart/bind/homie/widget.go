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
		var successMsg string

		switch r.URL().Query().Get("action") {
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

				if err == nil {
					successMsg = "Send signal of set config success"
				}
			}

		case "restart":
			err = bind.Restart(r.Context())
			if err == nil {
				successMsg = "Send restart signal success"
			}

		case "reset":
			err = bind.Reset(r.Context())
			if err == nil {
				successMsg = "Send reset signal success"
			}

		case "broadcast":
			err = r.Original().ParseForm()
			if err == nil {
				level := r.Original().PostFormValue("level")
				if level == "" {
					t.NotFound(w, r)
					return
				}

				err = bind.Broadcast(r.Context(), level, r.Original().PostFormValue("message"))
				if err == nil {
					successMsg = "Send broadcast message success"
				}
			}

		default:
			t.NotFound(w, r)
			return
		}

		if err != nil {
			w.SendJSON(response{
				Result:  "failed",
				Message: err.Error(),
			})

		} else {
			w.SendJSON(response{
				Result:  "success",
				Message: successMsg,
			})
		}

		return
	}

	vars := map[string]interface{}{
		"error":              err,
		"name":               "",
		"last_update":        bind.LastUpdate(),
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
