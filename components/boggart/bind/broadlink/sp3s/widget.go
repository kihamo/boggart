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
				widget.FlashError(r, "On failed with error %v", "", err)
			} else {
				widget.FlashSuccess(r, "On success", "")
			}
		} else {
			if err := b.Off(r.Context()); err != nil {
				widget.FlashError(r, "Off failed with error %v", "", err)
			} else {
				widget.FlashSuccess(r, "Off success", "")
			}
		}

		widget.Redirect(r.URL().Path, http.StatusFound, w, r)

		return
	}

	state, err := b.State()
	if err != nil {
		widget.FlashError(r, "Get state failed with error %v", "", err)
	}

	power, err := b.Power()
	if err != nil {
		widget.FlashError(r, "Get power failed with error %v", "", err)
	}

	widget.Render(r.Context(), "widget", map[string]interface{}{
		"state": state,
		"power": power,
	})
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
