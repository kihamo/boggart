package chromecast

import (
	"net/http"
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)

	data := map[string]interface{}{
		"url":    b.Config().(*Config).WidgetFileURL,
		"volume": bind.Volume(),
	}

	status := bind.status.Load()
	isPlaying := status == PlayerStatePlaying || status == PlayerStateBuffering

	if isPlaying {
		r.Session().FlashBag().Error(t.Translate(r.Context(), "Already playing", ""))
	}

	if r.IsPost() {
		data["volume"] = r.Original().FormValue("volume")
		data["url"] = r.Original().FormValue("url")

		if !isPlaying {
			volume, err := strconv.ParseInt(r.Original().FormValue("volume"), 10, 64)
			if err == nil {
				err = bind.SetVolume(r.Context(), volume)
				if err == nil {
					url := r.Original().FormValue("url")
					err = bind.PlayFromURL(r.Context(), url)
					data["url"] = url
				}
			}

			if err != nil {
				r.Session().FlashBag().Error(err.Error())
			} else {
				r.Session().FlashBag().Info(t.Translate(r.Context(), "File playing", ""))
				t.Redirect(r.URL().Path, http.StatusFound, w, r)
				return
			}
		}
	}

	t.Render(r.Context(), "widget", data)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
