package sun

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)

	t.Render(r.Context(), "widget", map[string]interface{}{
		"times": bind.Times(),
	})
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
