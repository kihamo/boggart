package handlers

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

type HikvisionHandler struct {
	dashboard.Handler

	Config config.Component
}

func (h *HikvisionHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	query := r.URL().Query()

	switch query.Get(":action") {
	case "preview":
		switch query.Get(":place") {
		case "street":
			if !h.Config.GetBool(boggart.ConfigHikvisionStreetEnabled) {
				return
			}

			isapi := hikvision.NewISAPI(
				h.Config.GetString(boggart.ConfigHikvisionStreetHost),
				h.Config.GetInt64(boggart.ConfigHikvisionStreetPort),
				h.Config.GetString(boggart.ConfigHikvisionStreetUsername),
				h.Config.GetString(boggart.ConfigHikvisionStreetPassword))

			image, err := isapi.StreamingPicture(h.Config.GetUint64Default(boggart.ConfigHikvisionStreetStreamingChannel, 101))
			if err != nil {
				h.NotFound(w, r)
				return
			}

			w.Header().Set("Content-Type", "image/jpeg; charset=\"UTF-8\"")
			w.Write(image)

			break

		default:
			h.NotFound(w, r)
		}

		break
	}
}
