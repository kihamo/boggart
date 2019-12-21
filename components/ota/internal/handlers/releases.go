package handlers

import (
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/ota"
	"github.com/kihamo/boggart/components/ota/release"
	"github.com/kihamo/boggart/components/ota/repository"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
)

type releaseView struct {
	ID           int
	Version      string
	Size         int64
	Checksum     string
	IsCurrent    bool
	Path         string
	Architecture string
}

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

type ReleasesHandler struct {
	dashboard.Handler

	Updater    *ota.Updater
	Repository *repository.MemoryRepository
}

func (h *ReleasesHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	releases, err := h.Repository.Releases("")
	if err != nil {
		r.Session().FlashBag().Error(err.Error())
	} else {
		switch q.Get(":action") {
		case "download":
			h.actionDownload(w, r, releases)
			return

		case "remove":
			h.actionRemove(w, r, releases)
			return

		case "upgrade":
			h.actionUpgrade(w, r, releases)
			return
		}
	}

	releasesView := make([]releaseView, 0, len(releases))
	for id, rl := range releases {
		rView := releaseView{
			ID:           id,
			Version:      rl.Version(),
			Size:         rl.Size(),
			Checksum:     hex.EncodeToString(rl.Checksum()),
			IsCurrent:    id == 0,
			Architecture: rl.Architecture(),
		}

		if releaseFile, ok := rl.(*release.LocalFileRelease); ok {
			rView.Path = releaseFile.Path()
		}

		releasesView = append(releasesView, rView)

	}

	h.Render(r.Context(), "releases", map[string]interface{}{
		"releases":    releasesView,
		"currentArch": runtime.GOARCH,
	})
}

func (h *ReleasesHandler) actionDownload(w *dashboard.Response, r *dashboard.Request, releases []ota.Release) {
	id, err := strconv.Atoi(r.URL().Query().Get(":id"))
	if err != nil {
		h.NotFound(w, r)
		return
	}

	for i, rl := range releases {
		if i == id {
			fileName := "release." + rl.Architecture() + ".bin"
			if releaseFile, ok := rl.(*release.LocalFileRelease); ok {
				fileName = filepath.Base(releaseFile.Path()) +
					"." + strings.ReplaceAll(releaseFile.Version(), " ", ".") +
					"." + rl.Architecture() + ".bin"
			}

			w.Header().Set("Content-Length", strconv.FormatInt(rl.Size(), 10))
			w.Header().Set("Content-Type", "application/x-binary")
			w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
			io.Copy(w, rl.BinFile())
			return
		}
	}

	h.NotFound(w, r)
	return
}

func (h *ReleasesHandler) actionRemove(w *dashboard.Response, r *dashboard.Request, releases []ota.Release) {
	if !r.IsPost() {
		h.MethodNotAllowed(w, r)
		return
	}

	id, err := strconv.Atoi(r.URL().Query().Get(":id"))
	if err != nil {
		h.NotFound(w, r)
		return
	}

	if id == 0 {
		err = errors.New("can't remove current release")
	} else {
		for i, rl := range releases {
			if i == id {
				h.Repository.Remove(rl)
				info := []interface{}{"version", rl.Version()}

				if releaseFile, ok := rl.(*release.LocalFileRelease); ok {
					os.Remove(releaseFile.Path())
					info = append(info, "path", releaseFile.Path())
				}

				logging.Log(r.Context()).Info("Remove release", info...)

				break
			}
		}
	}

	if err != nil {
		_ = w.SendJSON(response{
			Result:  "failed",
			Message: err.Error(),
		})
	} else {
		_ = w.SendJSON(response{
			Result: "success",
		})
	}
}

func (h *ReleasesHandler) actionUpgrade(w *dashboard.Response, r *dashboard.Request, releases []ota.Release) {
	if !r.IsPost() {
		h.MethodNotAllowed(w, r)
		return
	}

	id, err := strconv.Atoi(r.URL().Query().Get(":id"))
	if err != nil {
		h.NotFound(w, r)
		return
	}

	for i, rl := range releases {
		if i == id {
			err = h.Updater.Update(rl)
			if err != nil {
				r.Session().FlashBag().Error(err.Error())
			} else {
				info := []interface{}{"version", rl.Version()}

				if releaseFile, ok := rl.(*release.LocalFileRelease); ok {
					info = append(info, "path", releaseFile.Path())
				}

				logging.Log(r.Context()).Warn("Release upgrade", info...)
			}

			if r.URL().Query().Get("restart") != "" {
				err = h.Updater.Restart()
			}

			break
		}
	}

	if err != nil {
		_ = w.SendJSON(response{
			Result:  "failed",
			Message: err.Error(),
		})
	} else {
		_ = w.SendJSON(response{
			Result: "success",
		})
	}
}
