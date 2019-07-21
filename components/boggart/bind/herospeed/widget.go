package herospeed

import (
	"bytes"
	"io"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)

	query := r.URL().Query()
	action := query.Get("action")
	ctx := r.Context()
	cfg := b.Config().(*Config)

	vars := map[string]interface{}{
		"action":                   action,
		"preview_refresh_interval": cfg.PreviewRefreshInterval.Seconds(),
	}

	switch action {
	case "preview":
		buf := bytes.NewBuffer(nil)

		if err := bind.client.Snapshot(ctx, buf); err != nil {
			t.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Length", strconv.FormatInt(int64(buf.Len()), 10))

		if download := query.Get("download"); download != "" {
			filename := b.ID() + time.Now().Format("_20060102150405.jpg")

			w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
			w.Header().Set("Content-Type", "image/jpeg")
		}

		_, _ = io.Copy(w, buf)

		return

	case "configuration":
		configuration, err := bind.client.Configuration(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get configuration failed with error %s", "", err.Error()))
		}

		vars["configuration"] = configuration
	}

	t.Render(ctx, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
