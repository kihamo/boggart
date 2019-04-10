package hikvision

import (
	"bytes"
	"io"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/shadow/components/dashboard"
)

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)

	query := r.URL().Query()
	action := query.Get("action")

	vars := map[string]interface{}{
		"action":                   action,
		"preview_refresh_interval": b.Config().(*Config).PreviewRefreshInterval.Seconds(),
	}

	switch action {
	case "image":
		if r.IsPost() {
			var (
				ch  uint64
				err error
			)

			if channel := query.Get("channel"); channel == "" {
				ch = b.Config().(*Config).WidgetChannel
			} else {
				ch, err = strconv.ParseUint(channel, 10, 64)
				if err != nil {
					t.NotFound(w, r)
					return
				}
			}

			err = r.Original().ParseForm()
			if err == nil {
				for key, value := range r.Original().PostForm {
					if len(value) == 0 {
						continue
					}

					switch key {
					case "ir-cut-filter-type":
						err = bind.isapi.ImageIrCutFilter(r.Context(), ch, hikvision.ImageIrCutFilter{
							Type: hikvision.ImageIrCutFilterType(value[0]),
						})
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
					Result:  "success",
					Message: "Save success",
				})
			}

			return
		}

		response, err := bind.isapi.ImageChannels(r.Context())
		if err != nil {
			t.NotFound(w, r)
			return
		}

		vars["channels"] = response.Channels

	case "preview":
		var (
			ch  uint64
			err error
		)

		if channel := query.Get("channel"); channel == "" {
			ch = b.Config().(*Config).WidgetChannel
		} else {
			ch, err = strconv.ParseUint(channel, 10, 64)
			if err != nil {
				t.NotFound(w, r)
				return
			}
		}

		buf := bytes.NewBuffer(nil)
		if err = bind.Snapshot(r.Context(), ch, buf); err != nil {
			t.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Length", strconv.FormatInt(int64(buf.Len()), 10))

		if download := query.Get("download"); download != "" {
			filename := b.ID() + time.Now().Format("_20060102150405.jpg")

			w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
			w.Header().Set("Content-Type", "text/x-yaml")
		}

		_, _ = io.Copy(w, buf)

		return
	}

	t.Render(r.Context(), "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
