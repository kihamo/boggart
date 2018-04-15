package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

// easyjson:json
type deviceHandlerResponseSuccess struct {
	Result string `json:"result"`
}

// easyjson:json
type deviceHandlerResponseFailed struct {
	Result  string `json:"result"`
	Message string `json:"message"`
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
		var err error

		if device.IsEnabled() {
			err = device.Disable()
		} else {
			err = device.Enable()
		}

		if err == nil {
			w.SendJSON(deviceHandlerResponseSuccess{
				Result: "success",
			})
		} else {
			w.SendJSON(deviceHandlerResponseFailed{
				Result:  "failed",
				Message: err.Error(),
			})
		}
	default:
		h.NotFound(w, r)
		return
	}
}
