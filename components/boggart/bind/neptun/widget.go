package neptun

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	ctx := r.Context()
	provider := b.Provider()
	action := r.URL().Query().Get("action")

	vars := map[string]interface{}{
		"action": action,
	}

	switch action {
	case "counters":
		values, err := provider.CountersValues()
		if err != nil {
			widget.FlashError(r, "Get conters values failed with error %v", "", err)
		} else {
			vars["counter_values"] = values
		}

		configs, err := provider.CountersConfigurations()
		if err != nil {
			widget.FlashError(r, "Get counters configuration failed with error %v", "", err)
		} else {
			vars["counters_configs"] = configs
		}

	default:
		moduleCfg, err := provider.ModuleConfiguration()
		if err != nil {
			widget.FlashError(r, "Get module configuration failed with error %v", "", err)
		} else {
			vars["floor_washing"] = moduleCfg.FloorWashing()
			vars["first_group_alert"] = moduleCfg.FirstGroupAlert()
			vars["second_group_alert"] = moduleCfg.SecondGroupAlert()
			vars["wireless_sensor_low_battery"] = moduleCfg.WirelessSensorLowBattery()
			vars["wireless_sensor_loss"] = moduleCfg.WirelessSensorLoss()
			vars["first_group_tap_closing"] = moduleCfg.FirstGroupTapClosing()
			vars["second_group_tap_closing"] = moduleCfg.SecondGroupTapClosing()
			vars["wireless_sensor_pairing_mode"] = moduleCfg.WirelessSensorPairingMode()
			vars["first_group_tap_state"] = moduleCfg.FirstGroupTapState()
			vars["second_group_tap_state"] = moduleCfg.SecondGroupTapState()
			vars["two_groups_mode"] = moduleCfg.TwoGroupsMode()
			vars["taps_closing_on_sensor_loss"] = moduleCfg.TapsClosingOnSensorLoss()
			vars["keyboard_lock"] = moduleCfg.KeyboardLock()
		}
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
