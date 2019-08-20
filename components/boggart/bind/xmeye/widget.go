package xmeye

import (
	"io"
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

	vars := map[string]interface{}{
		"action": action,
	}

	switch action {
	case "logs":
		logs, err := bind.client.LogSearch(ctx, time.Now().Add(-time.Hour), time.Now(), 0)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get logs failed with error %s", "", err.Error()))
		}

		vars["logs"] = logs

	case "logs-export":
		reader, err := bind.client.LogExport(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Export logs failed with error %s", "", err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/zip")

		filename := b.ID() + time.Now().Format("_logs_20060102150405.zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

		_, _ = io.Copy(w, reader)

		return

	case "configs-export":
		reader, err := bind.client.ConfigExport(ctx)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Export config failed with error %s", "", err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/zip")

		filename := b.ID() + time.Now().Format("_config_20060102150405.zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

		_, _ = io.Copy(w, reader)

		return
	}

	t.Render(ctx, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
