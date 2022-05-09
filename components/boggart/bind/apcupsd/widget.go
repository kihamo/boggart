package apcupsd

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(_ *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	ctx := r.Context()
	vars := map[string]interface{}{}

	status, err := b.client.Status(ctx)
	if err != nil {
		widget.FlashError(r, "Get status failed with error %v", "", err)
	} else {
		vars["status"] = status
	}

	events, err := b.client.Events(ctx)
	if err != nil {
		widget.FlashError(r, "Get events failed with error %v", "", err)
	} else {
		vars["events"] = events
	}

	widget.Render(r.Context(), "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
