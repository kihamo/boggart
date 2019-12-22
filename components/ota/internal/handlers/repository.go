package handlers

import (
	"encoding/hex"
	"io"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/ota"
	"github.com/kihamo/boggart/components/ota/release"
	"github.com/kihamo/boggart/components/ota/repository"
	"github.com/kihamo/shadow/components/dashboard"
)

type RepositoryHandler struct {
	dashboard.Handler

	Repository *repository.DirectoryRepository
}

func (h *RepositoryHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	arch := strings.TrimSpace(r.URL().Query().Get("architecture"))

	releases, err := h.Repository.Releases(arch)
	if err != nil {
		h.InternalError(w, r, err)
		return
	}

	if id := r.URL().Query().Get(":id"); id != "" {
		for _, rl := range releases {
			if rlID := release.GenerateReleaseID(rl); rlID == id {
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

	if !r.Config().Bool(ota.ConfigRepositoryServerEnabled) {
		h.NotFound(w, r)
		return
	}

	type item struct {
		Architecture string `json:"architecture"`
		Checksum     string `json:"checksum"`
		Size         int64  `json:"size"`
		Version      string `json:"version"`
		File         string `json:"file"`
	}

	items := make([]item, 0, len(releases))
	for _, rl := range releases {
		fileURL := &url.URL{
			Scheme: "http",
			Host:   r.Original().Host,
			Path:   "/ota/repository/" + release.GenerateReleaseID(rl) + "/release.bin",
		}

		if r.Original().TLS != nil {
			fileURL.Scheme = "https"
		}

		items = append(items, item{
			Architecture: rl.Architecture(),
			Checksum:     hex.EncodeToString(rl.Checksum()),
			Size:         rl.Size(),
			Version:      rl.Version(),
			File:         fileURL.String(),
		})
	}

	_ = w.SendJSON(items)
}
