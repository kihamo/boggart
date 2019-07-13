package fcm

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	if r.IsPost() {
		ctx := r.Context()

		err := b.Bind().(*Bind).Send(ctx, r.Original().FormValue("message"))

		if err != nil {
			r.Session().FlashBag().Error(err.Error())
		} else {
			r.Session().FlashBag().Info(t.Translate(ctx, "Message sent", ""))
			t.Redirect(r.URL().Path, http.StatusFound, w, r)
			return
		}
	}

	t.Render(r.Context(), "widget", nil)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
