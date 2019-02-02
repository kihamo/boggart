package chromecast

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
	domain := dashboard.TemplateNamespaceFromContext(r.Context())

	status := bind.status.Load()
	isPlaying := status == "PLAYING" || status == "BUFFERING"
	if isPlaying {
		data["error"] = i18n.Locale(r.Context()).Translate(domain, "Already playing", "")
	}

	if r.IsPost() {
		data["volume"] = r.Original().FormValue("volume")
		data["url"] = r.Original().FormValue("url")

		if !isPlaying {
			var volume int64

			volume, err = strconv.ParseInt(r.Original().FormValue("volume"), 10, 64)
			if err == nil {
				bind.SetVolume(r.Context(), volume)

				url := r.Original().FormValue("url")
				err = bind.PlayFromURL(r.Context(), url)
				data["url"] = url
			}

			if err != nil {
				data["error"] = err.Error()
			} else {
				data["message"] = i18n.Locale(r.Context()).Translate(domain, "File playing", "")
			}
		}
	}

	t.Render(r.Context(), "widget", data)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}