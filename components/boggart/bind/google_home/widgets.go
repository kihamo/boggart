package google_home

import (
	"fmt"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	fmt.Println("33333")
	fmt.Println(r.URL().String())
	fmt.Println(b.Config())

	t.Render(r.Context(), "bind", nil)
}

func (t Type) WidgetTemplates() *assetfs.AssetFS {
	return assetFS()
}
