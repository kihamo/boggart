package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

// easyjson:json
type deviceHandlerResponseSuccess struct {
	Result string `json:"result"`
}

type DeviceHandler struct {
	dashboard.Handler
}

func (h *DeviceHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	query := r.URL().Query()

	deviceId := query.Get(":device")
	if deviceId == "" {
		h.NotFound(w, r)
		return
	}

	action := query.Get(":action")
	if action == "" {
		h.NotFound(w, r)
		return
	}

	device := r.Component().(boggart.Component).DevicesManager().Device(deviceId)
	if device == nil {
		h.NotFound(w, r)
		return
	}

	switch action {
	case "toggle":
		if device.IsEnabled() {
			device.Disable()
		} else {
			device.Enable()
		}

		w.SendJSON(deviceHandlerResponseSuccess{
			Result: "success",
		})
	default:
		h.NotFound(w, r)
		return
	}
}
