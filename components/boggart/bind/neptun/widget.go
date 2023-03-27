package neptun

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	var vars map[string]interface{}

	action := r.URL().Query().Get("action")

	switch action {
	default:
		action = "default"

		vars = b.widgetActionDefault(w, r)
	}

	vars["action"] = action
	b.Widget().RenderLayout(r.Context(), action, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}

func (b *Bind) widgetActionDefault(_ *dashboard.Response, r *dashboard.Request) map[string]interface{} {
	moduleCfg, err := b.Provider().ModuleConfiguration()
	if err != nil {
		b.Widget().FlashError(r, "Get module configuration failed with error %v", "", err)
	}

	return map[string]interface{}{
		"floor_washing":                moduleCfg.FloorWashing(),
		"first_group_alert":            moduleCfg.FirstGroupAlert(),
		"second_group_alert":           moduleCfg.SecondGroupAlert(),
		"wireless_sensor_low_battery":  moduleCfg.WirelessSensorLowBattery(),
		"wireless_sensor_loss":         moduleCfg.WirelessSensorLoss(),
		"first_group_tap_closing":      moduleCfg.FirstGroupTapClosing(),
		"second_group_tap_closing":     moduleCfg.SecondGroupTapClosing(),
		"wireless_sensor_pairing_mode": moduleCfg.WirelessSensorPairingMode(),
		"first_group_tap_state":        moduleCfg.FirstGroupTapState(),
		"second_group_tap_state":       moduleCfg.SecondGroupTapState(),
		"two_groups_mode":              moduleCfg.TwoGroupsMode(),
		"taps_closing_on_sensor_loss":  moduleCfg.TapsClosingOnSensorLoss(),
		"keyboard_lock":                moduleCfg.KeyboardLock(),
	}
}
