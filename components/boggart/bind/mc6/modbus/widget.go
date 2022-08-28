package modbus

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(_ *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	ctx := r.Context()
	provider := b.Provider()

	action := r.URL().Query().Get("action")

	deviceType, err := b.DeviceType(ctx)
	if err != nil {
		widget.FlashError(r, "Get device type failed with error %v", "", err)
	}

	vars := map[string]interface{}{
		"action": action,
		"device": deviceType,
	}

	switch action {
	case "sensors":
		if deviceType.IsSupportedTemperatureFormat() {
			value, err := provider.TemperatureFormat()
			vars["temperature_format"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedRoomTemperature() {
			value, err := provider.RoomTemperature()
			vars["room_temperature"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedFloorTemperature() {
			value, err := provider.FloorTemperature()
			vars["floor_temperature"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedFloorOverheat() {
			value, err := provider.FloorOverheat()
			vars["floor_overheat"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedHumidity() {
			value, err := provider.Humidity()
			vars["humidity"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

	case "status":
		if deviceType.IsSupportedStatus() {
			value, err := provider.Status()
			vars["power"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedHeatingValve() {
			value, err := provider.HeatingValve()
			vars["heating_valve"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedCoolingValve() {
			value, err := provider.CoolingValve()
			vars["cooling_valve"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedHeatingOutput() {
			value, err := provider.HeatingOutput()
			vars["heating_output"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedWindowsOpen() {
			value, err := provider.WindowsOpen()
			vars["windows_open"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

	case "away":
		if deviceType.IsSupportedAway() {
			value, err := provider.Away()
			vars["away"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedTemperatureFormat() {
			value, err := provider.TemperatureFormat()
			vars["temperature_format"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedAwayTemperature() {
			value, err := provider.AwayTemperature()
			vars["away_temperature"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

	case "hold":
		if deviceType.IsSupportedHoldingFunction() {
			value, err := provider.HoldingFunction()
			vars["holding_function"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedHoldingTimeHi() {
			value, err := provider.HoldingTimeHi()
			vars["holding_time_hi"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedHoldingTimeLow() {
			value, err := provider.HoldingTimeLow()
			vars["holding_time_low"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedTemperatureFormat() {
			value, err := provider.TemperatureFormat()
			vars["temperature_format"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedHoldingTemperature() {
			value, err := provider.HoldingTemperature()
			vars["holding_temperature"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

	case "fan":
		if deviceType.IsSupportedFanSpeedMode() {
			value, err := provider.FanSpeedMode()
			vars["fan_speed_mode"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

		if deviceType.IsSupportedFanSpeed() {
			value, err := provider.FanSpeed()
			vars["fan_speed"] = map[string]interface{}{
				"value": value,
				"error": err,
			}
		}

	default:
		if err == nil {
			vars["device_type"] = map[string]interface{}{
				"value": deviceType,
				"error": err,
			}

			if deviceType.IsSupportedSystemMode() {
				value, err := provider.SystemMode()
				vars["system_mode"] = map[string]interface{}{
					"value": value,
					"error": err,
				}
			}

			if deviceType.IsSupportedTemperatureFormat() {
				value, err := provider.TemperatureFormat()
				vars["temperature_format"] = map[string]interface{}{
					"value": value,
					"error": err,
				}
			}

			if deviceType.IsSupportedTargetTemperature() {
				value, err := provider.TargetTemperature()
				vars["target_temperature"] = map[string]interface{}{
					"value": value,
					"error": err,
				}
			}

			if deviceType.IsSupportedTargetTemperatureMaximum() {
				value, err := provider.TargetTemperatureMaximum()
				vars["target_temperature_maximum"] = map[string]interface{}{
					"value": value,
					"error": err,
				}
			}

			if deviceType.IsSupportedTargetTemperatureMinimum() {
				value, err := provider.TargetTemperatureMinimum()
				vars["target_temperature_minimum"] = map[string]interface{}{
					"value": value,
					"error": err,
				}
			}

			if deviceType.IsSupportedFloorTemperatureLimit() {
				value, err := provider.FloorTemperatureLimit()
				vars["floor_temperature_limit"] = map[string]interface{}{
					"value": value,
					"error": err,
				}
			}

			if deviceType.IsSupportedPanelLock() {
				value, err := provider.PanelLock()
				vars["panel_lock"] = map[string]interface{}{
					"value": value,
					"error": err,
				}
			}

			if deviceType.IsSupportedPanelLockPin1() {
				value, err := provider.PanelLockPin1()
				vars["panel_lock_pin_1"] = map[string]interface{}{
					"value": value,
					"error": err,
				}
			}

			if deviceType.IsSupportedPanelLockPin2() {
				value, err := provider.PanelLockPin2()
				vars["panel_lock_pin_2"] = map[string]interface{}{
					"value": value,
					"error": err,
				}
			}

			if deviceType.IsSupportedPanelLockPin3() {
				value, err := provider.PanelLockPin3()
				vars["panel_lock_pin_3"] = map[string]interface{}{
					"value": value,
					"error": err,
				}
			}

			if deviceType.IsSupportedPanelLockPin4() {
				value, err := provider.PanelLockPin4()
				vars["panel_lock_pin_4"] = map[string]interface{}{
					"value": value,
					"error": err,
				}
			}
		}
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
