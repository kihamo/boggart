package handlers

import (
	"encoding/hex"
	"errors"
	"net/http"
	"runtime"

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
	if r.IsPost() {
		file, header, err := r.Original().FormFile("release")

		if err == nil {
			defer file.Close()
			t := header.Header.Get("Content-Type")

			switch t {
			case "application/macbinary":
				var rl *release.LocalFileRelease

				rl, err = release.NewLocalFileFromStream(file, "", r.Config().String(ota.ConfigReleasesDirectory))
				if err == nil {
					h.Repository.Add(rl)

					_ = w.SendJSON(struct {
						ID           string `json:"id"`
						Version      string `json:"version"`
						Checksum     string `json:"checksum"`
						Architecture string `json:"architecture"`
						Size         int64  `json:"size"`
					}{
						ID:           release.GenerateReleaseID(rl),
						Version:      rl.Version(),
						Checksum:     hex.EncodeToString(rl.Checksum()),
						Architecture: rl.Architecture(),
						Size:         rl.Size(),
					})

					return
				}

			default:
				err = errors.New("unknown content type " + t)
			}
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}

		h.Redirect(r.Original().Referer(), http.StatusFound, w, r)
		return
	}

	h.Render(r.Context(), "update", map[string]interface{}{
		"goarch": runtime.GOARCH,
	})
}
