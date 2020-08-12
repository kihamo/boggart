package sun

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	b.Widget().Render(r.Context(), "widget", map[string]interface{}{
		"times": b.Times(),
	})
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
