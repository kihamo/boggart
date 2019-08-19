package xmeye

import (
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
		logs, err := bind.client.LogSearch(ctx, time.Now().Add(-time.Hour*24), time.Now(), 0)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get logs failed with error %s", "", err.Error()))
		}

		// change reverse
		for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
			logs[i], logs[j] = logs[j], logs[i]
		}

		vars["logs"] = logs
	}

	t.Render(ctx, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
