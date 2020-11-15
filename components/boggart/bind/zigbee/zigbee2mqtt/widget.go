package zigbee2mqtt

import (
	"fmt"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	ctx := r.Context()
	widget := b.Widget()

	if r.IsPost() {
		if r.URL().Query().Get("action") != "settings" {
			widget.NotFound(w, r)
			return
		}

		var (
			err        error
			successMsg string
		)

		err = r.Original().ParseForm()
		if err == nil {
			keys := make([]string, 0)

			for key, value := range r.Original().PostForm {
				if len(value) == 0 {
					continue
				}

				switch key {
				case "permit-join":
					err = b.SetPermitJoin(ctx, value[0] == "on")
				case "log-level":
					err = b.SetLogLevel(ctx, value[0])
				}

				if err != nil {
					err = fmt.Errorf("change setting %s return error: %w", key, err)
					break
				}

				keys = append(keys, key)
			}

			if err == nil {
				successMsg = "Change variables (" + strings.Join(keys, ",") + ") value success"
			}
		}

		if err != nil {
			_ = w.SendJSON(response{
				Result:  "failed",
				Message: err.Error(),
			})
		} else {
			_ = w.SendJSON(response{
				Result:  "success",
				Message: successMsg,
			})
		}

		return
	}

	vars := map[string]interface{}{
		"settings": b.Settings(),
		"devices":  b.Devices(),
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
