package handlers

import (
	"encoding/hex"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

type IndexHandler struct {
	dashboard.Handler

	Config config.Component
}

func (h *IndexHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	errors := []string{}
	vars := map[string]interface{}{
		"pulsarDeviceAddress":    "",
		"pulsarTemperatureIn":    0,
		"pulsarTemperatureOut":   0,
		"pulsarTemperatureDelta": 0,
	}

	connection := pulsar.NewConnection(
		h.Config.GetString(boggart.ConfigPulsarSerialAddress),
		h.Config.GetDuration(boggart.ConfigPulsarSerialTimeout))

	var (
		deviceAddress []byte
		err           error
	)

	deviceAddressConfig := h.Config.GetString(boggart.ConfigPulsarDeviceAddress)
	if deviceAddressConfig == "" {
		deviceAddress, err = connection.DeviceAddress()
	} else {
		deviceAddress, err = hex.DecodeString(deviceAddressConfig)
	}

	if err != nil {
		errors = append(errors, "Pulsar DeviceAddress failed: "+err.Error())
	}

	if len(deviceAddress) == 4 {
		vars["pulsarDeviceAddress"] = hex.EncodeToString(deviceAddress)
		device := pulsar.NewDevice(deviceAddress, connection)

		if value, err := device.TemperatureIn(); err == nil {
			vars["pulsarTemperatureIn"] = value
		} else {
			errors = append(errors, "Pulsar TemperatureIn failed: "+err.Error())
		}

		if value, err := device.TemperatureOut(); err == nil {
			vars["pulsarTemperatureOut"] = value
		} else {
			errors = append(errors, "Pulsar pulsarTemperatureOut failed: "+err.Error())
		}

		if value, err := device.TemperatureDelta(); err == nil {
			vars["pulsarTemperatureDelta"] = value
		} else {
			errors = append(errors, "Pulsar TemperatureDelta failed: "+err.Error())
		}
	}

	vars["errors"] = errors

	h.Render(r.Context(), boggart.ComponentName, "index", vars)
}
