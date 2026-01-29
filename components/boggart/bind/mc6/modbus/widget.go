package modbus

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"

	"github.com/kihamo/shadow/components/dashboard"

	"github.com/kihamo/boggart/providers/mc6"
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
		if deviceType.IsSupported(mc6.AddressTemperatureFormat) {
			if value, err := provider.TemperatureFormat(); err != nil {
				widget.FlashError(r, "Get temperature format failed with error %v", "", err)
			} else {
				vars["temperature_format"] = value
			}
		}

		registers, err := provider.ReadAsMap(mc6.AddressRoomTemperature, mc6.AddressHumidity)
		if err != nil {
			widget.FlashError(r, "Get sensors value failed with error %v", "", err)
		} else {
			if k := mc6.AddressRoomTemperature; deviceType.IsSupported(k) {
				vars["room_temperature"] = registers[k].Temperature()
			}

			if k := mc6.AddressFloorTemperature; deviceType.IsSupported(k) {
				vars["floor_temperature"] = registers[k].Temperature()
			}

			if k := mc6.AddressHumidity; deviceType.IsSupported(k) {
				vars["humidity"] = registers[k].Uint()
			}
		}

	case "status":
		registers, err := provider.ReadAsMap(mc6.AddressHeatingValve, mc6.AddressSystemError)
		if err != nil {
			widget.FlashError(r, "Get statuses failed with error %v", "", err)
		} else {
			if k := mc6.AddressHeatingValve; deviceType.IsSupported(k) {
				vars["heating_valve"] = registers[k].Bool()
			}

			if k := mc6.AddressCoolingValve; deviceType.IsSupported(k) {
				vars["cooling_valve"] = registers[k].Bool()
			}

			if k := mc6.AddressValve; deviceType.IsSupported(k) {
				vars["valve"] = registers[k].Bool()
			}

			if k := mc6.AddressFanHigh; deviceType.IsSupported(k) {
				vars["fan_high"] = registers[k].Bool()
			}

			if k := mc6.AddressFanMedium; deviceType.IsSupported(k) {
				vars["fan_medium"] = registers[k].Bool()
			}

			if k := mc6.AddressFanLow; deviceType.IsSupported(k) {
				vars["fan_low"] = registers[k].Bool()
			}

			if k := mc6.AddressHeatingOutput; deviceType.IsSupported(k) {
				vars["heating_output"] = registers[k].Bool()
			}

			if k := mc6.AddressHeat; deviceType.IsSupported(k) {
				vars["heat"] = registers[k].Bool()
			}

			if k := mc6.AddressHotWater; deviceType.IsSupported(k) {
				vars["hot_water"] = registers[k].Bool()
			}

			if k := mc6.AddressTouchLock; deviceType.IsSupported(k) {
				vars["touch_lock"] = registers[k].Bool()
			}

			if k := mc6.AddressWindowsOpen; deviceType.IsSupported(k) {
				vars["windows_open"] = registers[k].Bool()
			}

			if k := mc6.AddressHolidayFunction; deviceType.IsSupported(k) {
				vars["holiday_function"] = registers[k].Bool()
			}

			if k := mc6.AddressHoldingFunction; deviceType.IsSupported(k) {
				vars["holding_function"] = registers[k].Bool()
			}

			if k := mc6.AddressBoostFunction; deviceType.IsSupported(k) {
				vars["boost_function"] = registers[k].Bool()
			}

			if k := mc6.AddressFloorOverheat; deviceType.IsSupported(k) {
				vars["floor_overheat"] = registers[k].Bool()
			}

			if k := mc6.AddressAuxiliaryHeat; deviceType.IsSupported(k) {
				vars["auxiliary_heat"] = registers[k].Bool()
			}

			if k := mc6.AddressFanSpeedNumbers; deviceType.IsSupported(k) {
				vars["fan_speed_numbers"] = registers[k].Uint()
			}

			if k := mc6.AddressSystemError; deviceType.IsSupported(k) {
				vars["system_error"] = registers[k].Uint()
			}
		}

	case "away":

	case "hold":

	case "fan":

	default:

	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
