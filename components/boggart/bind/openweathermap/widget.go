package openweathermap

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()

	current, err := b.Current(r.Context())
	if err != nil {
		widget.FlashError(r, "Get current weather failed with error %v", "", err)
	}

	widget.Render(r.Context(), "widget", map[string]interface{}{
		"current": current,
	})
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
