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
		var (
			image []byte
			err   error
		)

		switch query.Get(":place") {
		case "hall":
			if !h.Config.Bool(boggart.ConfigHikvisionHallEnabled) {
				return
			}

			isapi := hikvision.NewISAPI(
				h.Config.String(boggart.ConfigHikvisionHallHost),
				h.Config.Int64(boggart.ConfigHikvisionHallPort),
				h.Config.String(boggart.ConfigHikvisionHallUsername),
				h.Config.String(boggart.ConfigHikvisionHallPassword))

			image, err = isapi.StreamingPicture(h.Config.Uint64(boggart.ConfigHikvisionHallStreamingChannel))
			break

		case "street":
			if !h.Config.Bool(boggart.ConfigHikvisionStreetEnabled) {
				return
			}

			isapi := hikvision.NewISAPI(
				h.Config.String(boggart.ConfigHikvisionStreetHost),
				h.Config.Int64(boggart.ConfigHikvisionStreetPort),
				h.Config.String(boggart.ConfigHikvisionStreetUsername),
				h.Config.String(boggart.ConfigHikvisionStreetPassword))

			image, err = isapi.StreamingPicture(h.Config.Uint64(boggart.ConfigHikvisionStreetStreamingChannel))
			break

		default:
			h.NotFound(w, r)
		}

		if err != nil {
			h.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg; charset=\"UTF-8\"")
		w.Write(image)

		break
	}
}
