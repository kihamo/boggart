package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

type ManagerUnregisterHandler struct {
	dashboard.Handler

	devicesManager boggart.DevicesManager
}

func NewManagerUnregisterHandler(devicesManager boggart.DevicesManager) *ManagerUnregisterHandler {
	return &ManagerUnregisterHandler{
		devicesManager: devicesManager,
	}
}

func (h *ManagerUnregisterHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	device := h.devicesManager.Device(r.URL().Query().Get(":id"))
	if device == nil {
		h.NotFound(w, r)
		return
	}

	err := h.devicesManager.Unregister(device.ID())

	type response struct {
		Result  string `json:"result"`
		Message string `json:"message,omitempty"`
	}

	if err != nil {
		w.SendJSON(response{
			Result:  "failed",
			Message: err.Error(),
		})
		return
	}

	w.SendJSON(response{
		Result: "success",
	})
}
