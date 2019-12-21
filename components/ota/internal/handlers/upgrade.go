package handlers

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/kihamo/boggart/components/ota"
	"github.com/kihamo/boggart/components/ota/release"
	"github.com/kihamo/boggart/components/ota/repository"
	"github.com/kihamo/shadow/components/dashboard"
)

type UpgradeHandler struct {
	dashboard.Handler

	Updater    *ota.Updater
	Repository *repository.MemoryRepository
}

func (h *UpgradeHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()
	refer := r.Original().Referer()

	switch q.Get(":action") {
	case "upload":
		file, header, err := r.Original().FormFile("release")

		if err == nil {
			defer file.Close()
			t := header.Header.Get("Content-Type")

			switch t {
			case "application/macbinary":
				buf := bytes.NewBuffer(nil)
				_, err = buf.ReadFrom(file)

				if err == nil {
					var rl *release.LocalFileRelease

					rl, err = release.NewLocalFileFromStream(file, "", r.Config().String(ota.ConfigReleasesDirectory))
					if err == nil {
						h.Repository.Add(rl)
					}
				}

			default:
				err = errors.New("unknown content type " + t)
			}
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		} else {
			_, _ = w.Write([]byte("success"))
		}

		h.Redirect(refer, http.StatusFound, w, r)
		return

	case "confirm":
		h.Redirect(refer, http.StatusFound, w, r)
		return

	case "restart":
		if err := h.Updater.Restart(); err != nil {
			r.Session().FlashBag().Error(err.Error())
		}

		h.Redirect(refer, http.StatusFound, w, r)
		return
	}

	h.Render(r.Context(), "update", map[string]interface{}{})
}
