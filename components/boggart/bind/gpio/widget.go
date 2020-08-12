package gpio

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()

	if r.IsPost() {
		if r.Original().FormValue("level") != "" {
			if err := b.High(r.Context()); err != nil {
				widget.FlashError(r, "Set high level failed with error %v", "", err)
			} else {
				widget.FlashSuccess(r, "Set high level success", "")
			}
		} else {
			if err := b.Low(r.Context()); err != nil {
				widget.FlashError(r, "Set low level failed with error %v", "", err)
			} else {
				widget.FlashSuccess(r, "Set low level success", "")
			}
		}

		widget.Redirect(r.URL().Path, http.StatusFound, w, r)

		return
	}

	widget.Render(r.Context(), "widget", map[string]interface{}{
		"level": b.Read(),
		"mode":  b.Mode(),
	})
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
