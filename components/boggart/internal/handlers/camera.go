package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

type CameraHandler struct {
	dashboard.Handler
}

func (h *CameraHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	query := r.URL().Query()

	switch query.Get(":action") {
	case "preview":
		var device boggart.Device

		switch query.Get(":place") {
		case "hall":
			device = r.Component().(boggart.Component).DevicesManager().Device(boggart.DeviceIdCameraHall.String())

		case "street":
			device = r.Component().(boggart.Component).DevicesManager().Device(boggart.DeviceIdCameraStreet.String())

		default:
			h.NotFound(w, r)
			return
		}

		if device == nil || !device.IsEnabled() {
			h.NotFound(w, r)
			return
		}

		camera, ok := device.(boggart.Camera)
		if !ok {
			h.NotFound(w, r)
			return
		}

		image, err := camera.Snapshot(r.Context())
		if err != nil {
			h.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg; charset=\"UTF-8\"")
		w.Write(image)

		break
	}
}
