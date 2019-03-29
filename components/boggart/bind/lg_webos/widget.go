package lg_webos

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	data := make(map[string]interface{})

	if bind.Status() != boggart.BindStatusOnline {
		r.Session().FlashBag().Error(t.Translate(r.Context(), "Device is offline", ""))
	}

	if r.IsPost() {
		toast := r.Original().FormValue("toast")
		data["toast"] = toast

		if toast != "" {
			if err := bind.Toast(toast); err != nil {
				r.Session().FlashBag().Error(err.Error())
			} else {
				r.Session().FlashBag().Success(t.Translate(r.Context(), "Message send success", ""))
				t.Redirect(r.URL().Path, http.StatusFound, w, r)
				return
			}
		} else {
			data["error"] = t.Translate(r.Context(), "Message is empty", "")
		}
	}

	t.Render(r.Context(), "widget", data)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
