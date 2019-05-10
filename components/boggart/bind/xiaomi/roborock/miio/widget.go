package miio

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	ctx := r.Context()
	action := r.URL().Query().Get("action")

	vars := map[string]interface{}{
		"action": action,
	}

	switch action {
	default:
		status, err := bind.device.Status(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get status failed with error %s", "", err.Error()))
		}

		vars["status"] = status

		info, err := bind.device.Info(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get info failed with error %s", "", err.Error()))
		}

		vars["info"] = info

		wifi, err := bind.device.WiFiStatus(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get WiFi status failed with error %s", "", err.Error()))
		}

		vars["wifi"] = wifi
	}

	t.Render(ctx, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
