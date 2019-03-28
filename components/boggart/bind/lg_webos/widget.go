package lg_webos

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	data := make(map[string]interface{})

	if bind.Status() != boggart.BindStatusOnline {
		data["error"] = i18n.Locale(r.Context()).Translate(boggart.ComponentName+"-bind-"+b.Type(), "Device is offline", "")
	}

	if r.IsPost() {
		toast := r.Original().FormValue("toast")
		data["toast"] = toast

		if toast != "" {
			if err := bind.Toast(toast); err != nil {
				data["error"] = err.Error()
			} else {
				data["message"] = i18n.Locale(r.Context()).Translate(boggart.ComponentName+"-bind-"+b.Type(), "Message send success", "")
			}
		} else {
			data["error"] = i18n.Locale(r.Context()).Translate(boggart.ComponentName+"-bind-"+b.Type(), "Message is empty", "")
		}
	}

	t.Render(r.Context(), "widget", data)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
