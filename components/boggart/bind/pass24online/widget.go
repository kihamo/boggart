package pass24online

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/providers/pass24online/client/feed"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	vars := make(map[string]interface{})

	if response, err := b.provider.Feed.GetFeed(feed.NewGetFeedParams(), nil); err != nil {
		widget.FlashError(r, "Get feed failed with error %v", "", err)
	} else {
		vars["feed"] = response.GetPayload().Body.Collection
	}

	widget.Render(r.Context(), "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
