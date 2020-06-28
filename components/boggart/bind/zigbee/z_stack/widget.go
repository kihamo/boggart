package z_stack

import (
	"encoding/hex"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	ctx := r.Context()

	vars := make(map[string]interface{})
	errors := make([]string, 0)

	client, err := bind.getClient(ctx)
	if err == nil {
		vars["devices"] = client.Devices()
		vars["permit_join"] = client.PermitJoinEnabled()

		if version, err := client.SysVersion(ctx); err == nil {
			vars["version"] = version
		} else {
			errors = append(errors, err.Error())
		}

		if info, err := client.UtilGetDeviceInfo(ctx); err == nil {
			vars["info"] = map[string]interface{}{
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
	} else {
		errors = append(errors, err.Error())
	}

	vars["errors"] = errors
	t.Render(ctx, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
