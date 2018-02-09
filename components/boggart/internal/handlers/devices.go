package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

type devicesHandlerDevice struct {
	Id          string
	Description string
	Types       []string
	Enabled     bool
}

type DevicesHandler struct {
	dashboard.Handler

	DeviceManager boggart.DeviceManager
}

func (h *DevicesHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	devices := h.DeviceManager.Devices()
	viewDevices := make([]*devicesHandlerDevice, 0, len(devices))

	for _, d := range devices {
		viewDevice := &devicesHandlerDevice{
			Id:          d.Id(),
			Description: d.Description(),
			Types:       []string{},
			Enabled:     d.IsEnabled(),
		}

		if _, ok := d.(boggart.Camera); ok {
			viewDevice.Types = append(viewDevice.Types, "camera")
		}

		viewDevices = append(viewDevices, viewDevice)
	}

	h.Render(r.Context(), boggart.ComponentName, "devices", map[string]interface{}{
		"devices": viewDevices,
	})
}
