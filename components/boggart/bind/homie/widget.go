package homie

import (
	"errors"
	"net/http"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
)

const (
	defaultOTATimeout = time.Minute
)

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

func (t Type) Widget(w *dashboard.Response, r *dashboard.Request, b boggart.BindItem) {
	bind := b.Bind().(*Bind)
	otaTimeout := defaultOTATimeout
	var err error

	if r.IsPost() {
		var successMsg string

		switch r.URL().Query().Get("action") {
		case "settings":
			err = r.Original().ParseForm()
			if err == nil {
				for key, value := range r.Original().PostForm {
					if len(value) == 0 {
						continue
					}

					err = bind.SettingsSet(r.Context(), key, value[0])
					if err != nil {
						break
					}
				}

				if err == nil {
					successMsg = "Send signal of set config success"
				}
			}

		case "restart":
			err = bind.Restart(r.Context())
			if err == nil {
				successMsg = "Send restart signal success"
			}

		case "reset":
			err = bind.Reset(r.Context())
			if err == nil {
				successMsg = "Send reset signal success"
			}

		case "broadcast":
			err = r.Original().ParseForm()
			if err == nil {
				level := r.Original().PostFormValue("level")
				if level == "" {
					t.NotFound(w, r)
					return
				}

				err = bind.Broadcast(r.Context(), level, r.Original().PostFormValue("message"))
				if err == nil {
					successMsg = "Send broadcast message success"
				}
			}

		case "ota":
			if v := r.Original().PostFormValue("timeout"); v != "" {
				var t time.Duration

				t, err = time.ParseDuration(v)
				if err == nil {
					otaTimeout = t
				}
			}

			if err == nil {
				file, header, err := r.Original().FormFile("firmware")

				if err == nil {
					defer file.Close()
					t := header.Header.Get("Content-Type")

					switch t {
					case "application/macbinary":
						err = bind.OTA(r.Context(), file, otaTimeout)

					default:
						err = errors.New("unknown content type " + t)
					}
				}
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			} else {
				w.Write([]byte("success"))
			}

			return

		default:
			t.NotFound(w, r)
			return
		}

		if err != nil {
			_ = w.SendJSON(response{
				Result:  "failed",
				Message: err.Error(),
			})

		} else {
			_ = w.SendJSON(response{
				Result:  "success",
				Message: successMsg,
			})
		}

		return
	}

	otaWritten, otaTotal := bind.OTAProgress()
	vars := map[string]interface{}{
		"error":              err,
		"name":               "",
		"last_update":        bind.LastUpdate(),
		"devices_attributes": bind.DeviceAttributes(),
		"nodes":              bind.Nodes(),
		"ota_running":        bind.OTAIsRunning(),
		"ota_written":        otaWritten,
		"ota_total":          otaTotal,
		"ota_checksum":       bind.OTAChecksum(),
		"ota_progress":       (float64(otaWritten) * float64(100)) / float64(otaTotal),
		"ota_timeout":        otaTimeout,
		"settings":           bind.SettingsAll(),
	}

	if attribute, ok := bind.DeviceAttribute("name"); ok {
		vars["name"] = attribute
	}

	t.Render(r.Context(), "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
