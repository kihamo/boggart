package ds18b20

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)

	value, err := bind.Temperature()

	t.Render(r.Context(), "widget", map[string]interface{}{
		"value": value,
		"error": err,
	})
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
