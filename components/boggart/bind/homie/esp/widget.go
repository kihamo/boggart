package esp

import (
	"errors"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
	tm "github.com/kihamo/shadow/misc/time"
)

const (
	defaultOTATimeout = time.Minute
)

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	otaTimeout := defaultOTATimeout
	widget := b.Widget()

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

					err = b.settingsSet(r.Context(), key, value[0])
					if err != nil {
						break
					}
				}

				if err == nil {
					successMsg = "Send signal of set config success"
				}
			}

		case "restart":
			err = b.Restart(r.Context())
			if err == nil {
				successMsg = "Send restart signal success"
			}

		case "reset":
			err = b.Reset(r.Context())
			if err == nil {
				successMsg = "Send reset signal success"
			}

		case "broadcast":
			err = r.Original().ParseForm()
			if err == nil {
				level := r.Original().PostFormValue("level")
				if level == "" {
					widget.NotFound(w, r)
					return
				}

				err = b.Broadcast(r.Context(), level, r.Original().PostFormValue("message"))
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
				var (
					file   multipart.File
					header *multipart.FileHeader
				)

				file, header, err = r.Original().FormFile("firmware")

				if err == nil {
					defer file.Close()

					t := header.Header.Get("Content-Type")

					switch t {
					case "application/macbinary":
						err = b.OTA(r.Context(), file, otaTimeout)

					default:
						err = errors.New("unknown content type " + t)
					}
				}
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
			} else {
				_, _ = w.Write([]byte("success"))
			}

			return

		default:
			widget.NotFound(w, r)
			return
		}

		if err != nil {
			_ = w.SendJSON(boggart.NewResponseJSON().FailedError(err))
		} else {
			_ = w.SendJSON(boggart.NewResponseJSON().Success(successMsg))
		}

		return
	}

	otaWritten, otaTotal := b.OTAProgress()

	name, _ := b.DeviceAttribute("name")

	protocol, ok := b.DeviceAttribute("homie")
	if ok && len(protocol.(string)) > 0 {
		protocol = protocol.(string)[:1]
	}

	lastUpdate := b.LastUpdate()
	vars := map[string]interface{}{
		"error":              err,
		"name":               name,
		"protocol_major":     protocol,
		"last_update":        lastUpdate,
		"devices_attributes": b.DeviceAttributes(),
		"nodes":              b.nodesList(),
		"ota_enabled":        b.OTAIsEnabled(),
		"ota_running":        b.OTAIsRunning(),
		"ota_written":        otaWritten,
		"ota_total":          otaTotal,
		"ota_checksum":       b.OTAChecksum(),
		"ota_progress":       (float64(otaWritten) * float64(100)) / float64(otaTotal),
		"ota_timeout":        otaTimeout,
		"settings":           b.settingsAll(),
	}

	if lastUpdate != nil {
		vars["last_update_delta"] = tm.DateSinceAsMessage(*lastUpdate)
	}

	widget.Render(r.Context(), "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
