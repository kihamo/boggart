package zstack

import (
	"encoding/hex"
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
					if value[0] == "on" {
						err = b.client.PermitJoin(ctx, b.permitJoinDuration())
					} else {
						err = b.client.PermitJoinDisable(ctx)
					}

				case "led":
					err = b.client.LED(ctx, value[0] == "on")
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

	vars := make(map[string]interface{})
	errors := make([]string, 0)

	client, err := b.getClient()
	if err == nil {
		vars["devices"] = client.Devices()
		vars["led_support"] = client.LEDSupport(ctx)
		vars["led_enabled"] = client.LEDEnabled()
		vars["permit_join"] = client.PermitJoinEnabled()

		if version, err := client.Version(ctx); err == nil {
			vars["version"] = version
		} else {
			errors = append(errors, err.Error())
		}

		if info, err := client.UtilGetDeviceInfo(ctx); err == nil {
			vars["info"] = map[string]interface{}{
				"device_type":  info.DeviceType,
				"device_state": info.DeviceState,
				"ieee_address": hex.EncodeToString(info.IEEEAddr),
			}
		} else {
			errors = append(errors, err.Error())
		}

		if info, err := client.ZDOExtNetworkInfo(ctx); err == nil {
			vars["network"] = map[string]interface{}{
				"pan_id":          info.PanID,
				"extended_pan_id": hex.EncodeToString(info.ExtendedPanID),
				"channel":         info.Channel,
			}
		} else {
			errors = append(errors, err.Error())
		}

		if check, err := client.InitializationCheck(ctx); err == nil {
			vars["checks"] = map[string]interface{}{
				"init": check,
			}
		} else {
			errors = append(errors, err.Error())
		}
	} else {
		errors = append(errors, err.Error())
	}

	vars["errors"] = errors
	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
