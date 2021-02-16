package herospeed

import (
	"bytes"
	"io"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	query := r.URL().Query()
	action := query.Get("action")
	ctx := r.Context()
	widget := b.Widget()

	vars := map[string]interface{}{
		"action":                   action,
		"preview_refresh_interval": b.config().PreviewRefreshInterval.Seconds(),
	}

	switch action {
	case "preview":
		buf := bytes.NewBuffer(nil)

		if err := b.client.Snapshot(ctx, buf); err != nil {
			widget.InternalError(w, r, err)
			return
		}

		w.Header().Set("Content-Length", strconv.FormatInt(int64(buf.Len()), 10))

		if download := query.Get("download"); download != "" {
			filename := b.Meta().ID() + time.Now().Format("_20060102150405.jpg")

			w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
			w.Header().Set("Content-Type", "image/jpeg")
		}

		_, _ = io.Copy(w, buf)

		return

	case "configuration":
		configuration, err := b.client.Configuration(ctx)
		if err != nil {
			widget.FlashError(r, "Get configuration failed with error %v", "", err)
		}

		vars["configuration"] = configuration
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
