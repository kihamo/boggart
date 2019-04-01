package gpio

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)

	if r.IsPost() {
		if r.Original().FormValue("state") != "" {
			if err := bind.High(r.Context()); err != nil {
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Set high failed with error %s", "", err.Error()))
			} else {
				r.Session().FlashBag().Success(t.Translate(r.Context(), "Set high success", ""))
			}
		} else {
			if err := bind.Low(r.Context()); err != nil {
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Set low failed with error %s", "", err.Error()))
			} else {
				r.Session().FlashBag().Success(t.Translate(r.Context(), "Set low success", ""))
			}
		}

		t.Redirect(r.URL().Path, http.StatusFound, w, r)
		return
	}

	t.Render(r.Context(), "widget", map[string]interface{}{
		"state": bind.Read(),
	})
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
