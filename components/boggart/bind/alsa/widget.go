package alsa

import (
	"net/http"
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/shadow/components/dashboard"
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	widget := b.Widget()
	var err error

	data := map[string]interface{}{
		"url":    b.config.WidgetFileURL,
		"volume": b.Volume(),
	}

	isPlaying := b.PlayerStatus() == StatusPlaying
	if isPlaying {
		widget.FlashError(r, "Already playing", "")
	}

	if r.IsPost() {
		data["volume"] = r.Original().FormValue("volume")
		data["url"] = r.Original().FormValue("url")

		if !isPlaying {
			var volume int64

			volume, err = strconv.ParseInt(r.Original().FormValue("volume"), 10, 64)
			if err == nil {
				err = b.SetVolume(volume)
				if err == nil {
					url := r.Original().FormValue("url")
					err = b.PlayFromURL(url)
					data["url"] = url
				}
			}

			if err != nil {
				widget.FlashError(r, err.Error(), "")
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
