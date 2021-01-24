package ds18b20

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()

	values, err := b.Temperatures()
	if err != nil {
		widget.FlashError(r, "Get values failed %v", "", err)
	}

	b.Widget().Render(r.Context(), "widget", map[string]interface{}{
		"values": values,
	})
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
