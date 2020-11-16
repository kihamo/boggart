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
	vars := map[string]interface{}{}

	action := r.URL().Query().Get("action")
	vars["action"] = action

	switch action {
	case "settings":
		if !r.IsPost() {
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

	case "map":
		network, err := b.NetworkMap(ctx)
		if err != nil {
			widget.FlashError(r, "Get network map failed with error %v", "", err)
		} else {
			vars["network_map"] = network

			for _, node := range network.Nodes {
				if node.Type == "Coordinator" {
					vars["coordinator_address"] = node.IEEEAddress
					break
				}
			}
		}

	default:
		vars["settings"] = b.Settings()
		vars["devices"] = b.Devices()
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
