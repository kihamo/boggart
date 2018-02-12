package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/shadow/components/dashboard"
)

type devicesHandlerDevice struct {
	TasksCount  int
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
	list := h.DeviceManager.Devices()
	viewList := make([]*devicesHandlerDevice, 0, len(list))

	for _, d := range list {
		viewDevice := &devicesHandlerDevice{
			Id:          d.Id(),
			Description: d.Description(),
			Types:       []string{},
			TasksCount:  len(d.Tasks()),
			Enabled:     d.IsEnabled(),
		}

		if _, ok := d.(boggart.Camera); ok {
			viewDevice.Types = append(viewDevice.Types, "camera")
		}

		if _, ok := d.(boggart.Phone); ok {
			viewDevice.Types = append(viewDevice.Types, "phone")
		}

		if _, ok := d.(*devices.VideoRecorderHikVision); ok {
			viewDevice.Types = append(viewDevice.Types, "video recorder")
		}

		viewList = append(viewList, viewDevice)
	}

	h.Render(r.Context(), boggart.ComponentName, "devices", map[string]interface{}{
		"devices": viewList,
	})
}
