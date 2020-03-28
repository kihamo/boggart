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
		if r.Original().FormValue("level") != "" {
			if err := bind.High(r.Context()); err != nil {
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Set high level failed with error %s", "", err.Error()))
			} else {
				r.Session().FlashBag().Success(t.Translate(r.Context(), "Set high level success", ""))
			}
		} else {
			if err := bind.Low(r.Context()); err != nil {
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Set low level failed with error %s", "", err.Error()))
			} else {
				r.Session().FlashBag().Success(t.Translate(r.Context(), "Set low level success", ""))
			}
		}

		t.Redirect(r.URL().Path, http.StatusFound, w, r)

		return
	}

	t.Render(r.Context(), "widget", map[string]interface{}{
		"level": bind.Read(),
		"mode":  bind.Mode(),
	})
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
