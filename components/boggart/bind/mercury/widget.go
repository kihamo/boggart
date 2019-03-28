package mercury

import (
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	vars := map[string]interface{}{}

	// date time
	now := time.Now()
	v, err := bind.provider.Datetime()
	variable := map[string]interface{}{
		"value": v,
		"delta": 0,
		"error": err,
	}
	if err == nil {
		variable["delta"] = int64(now.Sub(v).Seconds())
	}

	vars["datetime"] = variable

	// make date
	v, err = bind.provider.MakeDate()
	vars["make_date"] = map[string]interface{}{
		"value": v,
		"error": err,
	}

	// last power off
	v, err = bind.provider.LastPowerOffDatetime()
	vars["last_power_off_datetime"] = map[string]interface{}{
		"value": v,
		"error": err,
	}

	// last power on
	v, err = bind.provider.LastPowerOnDatetime()
	vars["last_power_on_datetime"] = map[string]interface{}{
		"value": v,
		"error": err,
	}

	t.Render(r.Context(), "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
