package sp3s

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()

	if r.IsPost() {
		if r.Original().FormValue("state") != "" {
			if err := b.On(r.Context()); err != nil {
				r.Session().FlashBag().Error(widget.Translate(r.Context(), "On failed with error %s", "", err.Error()))
			} else {
				r.Session().FlashBag().Success(widget.Translate(r.Context(), "On success", ""))
			}
		} else {
			if err := b.Off(r.Context()); err != nil {
				r.Session().FlashBag().Error(widget.Translate(r.Context(), "Off failed with error %s", "", err.Error()))
			} else {
				r.Session().FlashBag().Success(widget.Translate(r.Context(), "Off success", ""))
			}
		}

		widget.Redirect(r.URL().Path, http.StatusFound, w, r)

		return
	}

	state, err := b.State()
	if err != nil {
		r.Session().FlashBag().Error(widget.Translate(r.Context(), "Get state failed with error %s", "", err.Error()))
	}

	power, err := b.Power()
	if err != nil {
		r.Session().FlashBag().Error(widget.Translate(r.Context(), "Get power failed with error %s", "", err.Error()))
	}

	widget.Render(r.Context(), "widget", map[string]interface{}{
		"state": state,
		"power": power,
	})
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
