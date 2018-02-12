package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
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
			Types:       make([]string, 0, len(d.Types())),
			TasksCount:  len(d.Tasks()),
			Enabled:     d.IsEnabled(),
		}

		for _, t := range d.Types() {
			viewDevice.Types = append(viewDevice.Types, t.String())
		}

		viewList = append(viewList, viewDevice)
	}

	h.Render(r.Context(), boggart.ComponentName, "devices", map[string]interface{}{
		"devices": viewList,
	})
}
