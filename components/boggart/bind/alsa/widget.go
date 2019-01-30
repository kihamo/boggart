package alsa

import (
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
)

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	var err error

	data := map[string]interface{}{
		"url":    b.Config().(*Config).WidgetFileURL,
		"volume": bind.Volume(),
	}

	isPlaying := bind.PlayerStatus() == StatusPlaying
	if isPlaying {
		data["error"] = i18n.Locale(r.Context()).Translate(boggart.ComponentName+"-bind-"+b.Type(), "Already playing", "")
	}

	if r.IsPost() {
		data["volume"] = r.Original().FormValue("volume")
		data["url"] = r.Original().FormValue("url")

		if !isPlaying {
			var volume int64

			volume, err = strconv.ParseInt(r.Original().FormValue("volume"), 10, 64)
			if err == nil {
				bind.SetVolume(volume)

				url := r.Original().FormValue("url")
				err = bind.PlayFromURL(url)
				data["url"] = url
			}

			if err != nil {
				data["error"] = err.Error()
			} else {
				data["message"] = i18n.Locale(r.Context()).Translate(boggart.ComponentName+"-bind-"+b.Type(), "File playing", "")
			}
		}
	}

	t.Render(r.Context(), "widget", data)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
