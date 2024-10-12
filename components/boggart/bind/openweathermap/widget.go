package openweathermap

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/providers/openweathermap"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	ctx := r.Context()

	b.locationMutex.RLock()
	vars := map[string]interface{}{
		"icon":          openweathermap.Icon,
		"location_name": b.locationName,
	}
	b.locationMutex.RUnlock()

	response, err := b.OneCallMigrate(ctx)
	if err != nil {
		widget.FlashError(r, "One call failed with error %v", "", err)
	} else {
		vars["current"] = response.Current
		vars["daily"] = response.Daily
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
