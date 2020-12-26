package wol

import (
	"net"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()

	if r.IsPost() {
		hw, err := net.ParseMAC(r.Original().FormValue("mac"))
		if err == nil {
			err = b.WOL(hw, nil, nil)
		}

		if err != nil {
			widget.FlashError(r, err, "")
		} else {
			widget.FlashInfo(r, "Sent magic packet to %s", "", hw.String())
			widget.Redirect(r.URL().Path, http.StatusFound, w, r)
			return
		}
	}

	widget.Render(r.Context(), "widget", nil)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
