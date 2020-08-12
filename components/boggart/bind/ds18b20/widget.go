package ds18b20

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	value, err := b.Temperature()

	b.Widget().Render(r.Context(), "widget", map[string]interface{}{
		"value": value,
		"error": err,
	})
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
