package handlers

import (
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
	Repository *repository.DirectoryRepository
}

func (h *UpgradeHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()
	refer := r.Original().Referer()

	switch q.Get(":action") {
	case "confirm":
		h.Redirect(refer, http.StatusFound, w, r)
		return

	case "restart":
		if err := h.Updater.Restart(); err != nil {
			r.Session().FlashBag().Error(err.Error())
		}

		h.Redirect(refer, http.StatusFound, w, r)
		return

	default:
		if r.IsPost() {
			var id int64
			file, header, err := r.Original().FormFile("release")

			if err == nil {
				defer file.Close()
				t := header.Header.Get("Content-Type")

				switch t {
				case "application/macbinary":
					var rl *release.LocalFileRelease

					rl, err = release.NewLocalFileFromStream(file, "", r.Config().String(ota.ConfigReleasesDirectory))
					if err == nil {
						id = h.Repository.Add(rl)
					}

				default:
					err = errors.New("unknown content type " + t)
				}
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
			} else {
				_ = w.SendJSON(struct {
					ID int64 `json:"id"`
				}{
					ID: id,
				})
			}

			h.Redirect(refer, http.StatusFound, w, r)
			return
		}
	}

	h.Render(r.Context(), "update", map[string]interface{}{})
}
