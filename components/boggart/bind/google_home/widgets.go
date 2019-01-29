package google_home

import (
	"fmt"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	fmt.Fprint(w, "widget")
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
