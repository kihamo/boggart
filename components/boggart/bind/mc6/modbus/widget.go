package modbus

import (
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(_ *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	vars := map[string]interface{}{}
	ctx := r.Context()
	provider := b.Provider()

	action := r.URL().Query().Get("action")
	vars["action"] = action

	deviceType, err := b.DeviceType(ctx)
	if err != nil {
		widget.FlashError(r, "Get device type failed with error %v", "", err)
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

		if deviceType.IsSupportedHoldingFunction() {
			value, err := provider.HoldingFunction()
			vars["holding_function"] = map[string]interface{}{
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

			if deviceType.IsSupportedFanSpeed() {
				value, err := provider.FanSpeed()
				vars["fan_speed"] = map[string]interface{}{
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
		}
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
