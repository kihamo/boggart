package chromecast

import (
	"net/http"
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()

	data := map[string]interface{}{
		"url":    b.config.WidgetFileURL,
		"volume": b.Volume(),
	}

	status := b.status.Load()
	isPlaying := status == PlayerStatePlaying || status == PlayerStateBuffering

	if isPlaying {
		widget.FlashError(r, "Already playing", "")
	}

	if r.IsPost() {
		data["volume"] = r.Original().FormValue("volume")
		data["url"] = r.Original().FormValue("url")

		if !isPlaying {
			volume, err := strconv.ParseInt(r.Original().FormValue("volume"), 10, 64)
			if err == nil {
				err = b.SetVolume(r.Context(), volume)
				if err == nil {
					url := r.Original().FormValue("url")
					err = b.PlayFromURL(r.Context(), url)
					data["url"] = url
				}
			}

			if err != nil {
				widget.FlashError(r, err, "")
			} else {
				widget.FlashInfo(r, "File playing", "")
				widget.Redirect(r.URL().Path, http.StatusFound, w, r)
				return
			}
		}
	}

	widget.Render(r.Context(), "widget", data)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
