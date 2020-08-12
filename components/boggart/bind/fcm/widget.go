package fcm

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()

	if r.IsPost() {
		ctx := r.Context()

		err := b.Send(ctx, r.Original().FormValue("message"))

		if err != nil {
			widget.FlashError(r, err.Error(), "")
		} else {
			widget.FlashInfo(r, "Message sent", "")
			widget.Redirect(r.URL().Path, http.StatusFound, w, r)
			return
		}
	}

	widget.Render(r.Context(), "widget", nil)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
