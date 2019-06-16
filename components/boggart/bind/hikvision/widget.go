package hikvision

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision2/client/image"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision2/client/streaming"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision2/client/system"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision2/models"
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
	ctx := r.Context()

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
						params := image.NewSetImageIrCutFilterParamsWithContext(ctx).
							WithChannel(ch).
							WithIrcutFilter(&models.IrcutFilter{
								Type: value[0],
							})

						_, err = bind.client.Image.SetImageIrCutFilter(params, nil)
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

		response, err := bind.client.Image.GetImageChannels(image.NewGetImageChannelsParamsWithContext(ctx), nil)
		if err != nil {
			t.NotFound(w, r)
			return
		}

		vars["channels"] = response.Payload

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
		params := streaming.NewGetStreamingPictureParamsWithContext(ctx).WithChannel(ch)
		if _, err := bind.client.Streaming.GetStreamingPicture(params, nil, buf); err != nil {
			t.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Length", strconv.FormatInt(int64(buf.Len()), 10))

		if download := query.Get("download"); download != "" {
			filename := b.ID() + time.Now().Format("_20060102150405.jpg")

			w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
			// w.Header().Set("Content-Type", "image/jpeg")
		}

		_, _ = io.Copy(w, buf)

		return

	case "system":
		if r.IsPost() {
			file, header, err := r.Original().FormFile("firmware")

			if err == nil {
				defer file.Close()
				t := header.Header.Get("Content-Type")

				switch t {
				case "application/octet-stream":
					bind.FirmwareUpdate(file)

				default:
					err = errors.New("unknown content type " + t)
				}
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
			} else {
				_, _ = w.Write([]byte("success"))
			}

			return
		}

		info, err := bind.client.System.GetSystemDeviceInfo(system.NewGetSystemDeviceInfoParamsWithContext(ctx), nil)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get device info failed with error %s", "", err.Error()))
		}

		vars["info"] = info.Payload

		upgrade, err := bind.client.System.GetSystemUpgradeStatus(system.NewGetSystemUpgradeStatusParamsWithContext(ctx), nil)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get upgrade status failed with error %s", "", err.Error()))
		}

		vars["upgrade"] = upgrade.Payload
	}

	t.Render(ctx, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
