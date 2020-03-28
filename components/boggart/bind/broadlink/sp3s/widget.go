package sp3s

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
			if err := bind.On(r.Context()); err != nil {
				r.Session().FlashBag().Error(t.Translate(r.Context(), "On failed with error %s", "", err.Error()))
			} else {
				r.Session().FlashBag().Success(t.Translate(r.Context(), "On success", ""))
			}
		} else {
			if err := bind.Off(r.Context()); err != nil {
				r.Session().FlashBag().Error(t.Translate(r.Context(), "Off failed with error %s", "", err.Error()))
			} else {
				r.Session().FlashBag().Success(t.Translate(r.Context(), "Off success", ""))
			}
		}

		t.Redirect(r.URL().Path, http.StatusFound, w, r)

		return
	}

	state, err := bind.State()
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(r.Context(), "Get state failed with error %s", "", err.Error()))
	}

	power, err := bind.Power()
	if err != nil {
		r.Session().FlashBag().Error(t.Translate(r.Context(), "Get power failed with error %s", "", err.Error()))
	}

	t.Render(r.Context(), "widget", map[string]interface{}{
		"state": state,
		"power": power,
	})
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
