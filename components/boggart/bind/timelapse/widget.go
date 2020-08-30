package timelapse

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	b.Capture(r.Context(), w)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return nil
}
