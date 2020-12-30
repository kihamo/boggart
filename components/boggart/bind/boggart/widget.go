package boggart

import (
	"net/http"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()

	if r.IsPost() {
		time.AfterFunc(time.Second*5, func() {
			b.application.Shutdown()
		})

		widget.FlashSuccess(r, "Will shutdown after 5 seconds", "")
		widget.Redirect(r.URL().Path, http.StatusFound, w, r)

		return
	}

	widget.Render(r.Context(), "widget", nil)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
