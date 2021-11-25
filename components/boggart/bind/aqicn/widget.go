package aqicn

import (
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	ctx := r.Context()
	now := time.Now()
	vars := map[string]interface{}{
		"current_day": time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
	}

	feed, err := b.Feed(ctx)
	if err != nil {
		widget.FlashError(r, "Get feed failed with error %v", "", err)
	} else {
		vars["feed"] = feed
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
