package hikvision

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/hikvision/client/event"
	"github.com/kihamo/boggart/providers/hikvision/client/image"
	"github.com/kihamo/boggart/providers/hikvision/client/streaming"
	"github.com/kihamo/boggart/providers/hikvision/client/system"
	"github.com/kihamo/boggart/providers/hikvision/models"
	static "github.com/kihamo/boggart/providers/hikvision/static/models"
	"github.com/kihamo/shadow/components/dashboard"
)

type response struct {
	Result  string `json:"result"`
	Message string `json:"message,omitempty"`
}

func (b *Bind) WidgetHandler(w *dashboard.Response, r *dashboard.Request) {
	query := r.URL().Query()
	action := query.Get("action")
	ctx := r.Context()
	widget := b.Widget()

	vars := map[string]interface{}{
		"action":                   action,
		"events_enabled":           b.config.EventsEnabled,
		"preview_refresh_interval": b.config.PreviewRefreshInterval.Seconds(),
	}

	switch action {
	case "image":
		if r.IsPost() {
			var (
				ch  uint64
				err error
			)

			if channel := query.Get("channel"); channel == "" {
				ch = b.config.WidgetChannel
			} else {
				ch, err = strconv.ParseUint(channel, 10, 64)
				if err != nil {
					widget.NotFound(w, r)
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

						_, err = b.client.Image.SetImageIrCutFilter(params, nil)

					case "flip":
						flip := &models.ImageFlip{
							Enabled: value[0] != "disabled",
						}

						if flip.Enabled {
							flip.Style = value[0]
						}

						params := image.NewSetImageFlipParamsWithContext(ctx).
							WithChannel(ch).
							WithImageFlip(flip)

						_, err = b.client.Image.SetImageFlip(params, nil)
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

		response, err := b.client.Image.GetImageChannels(image.NewGetImageChannelsParamsWithContext(ctx), nil)
		if err != nil {
			widget.NotFound(w, r)
			return
		}

		vars["channels"] = response.Payload

	case "preview":
		var (
			ch  uint64
			err error
		)

		if channel := query.Get("channel"); channel == "" {
			ch = b.config.WidgetChannel
		} else {
			ch, err = strconv.ParseUint(channel, 10, 64)
			if err != nil {
				widget.NotFound(w, r)
				return
			}
		}

		buf := bytes.NewBuffer(nil)
		params := streaming.NewGetStreamingPictureParamsWithContext(ctx).WithChannel(ch)

		if _, err := b.client.Streaming.GetStreamingPicture(params, nil, buf); err != nil {
			widget.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Length", strconv.FormatInt(int64(buf.Len()), 10))

		if download := query.Get("download"); download != "" {
			filename := b.Meta().ID() + time.Now().Format("_20060102150405.jpg")

			w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
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
					b.FirmwareUpdate(file)

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

		info, err := b.client.System.GetSystemDeviceInfo(system.NewGetSystemDeviceInfoParamsWithContext(ctx), nil)
		if err != nil {
			widget.FlashError(r, "Get device info failed with error %v", "", err)
		}

		vars["info"] = info.Payload

		upgrade, err := b.client.System.GetSystemUpgradeStatus(system.NewGetSystemUpgradeStatusParamsWithContext(ctx), nil)
		if err != nil {
			widget.FlashError(r, "Get upgrade status failed with error %v", "", err)
		}

		vars["upgrade"] = upgrade.Payload

	case "notification":
		if !b.config.EventsEnabled {
			widget.NotFound(w, r)
			return
		}

		params := event.NewGetNotificationHTTPHostParamsWithContext(ctx).
			WithHTTPHost(1)

		notification, err := b.client.Event.GetNotificationHTTPHost(params, nil)
		if err != nil {
			widget.FlashError(r, "Get notification http host failed with error %v", "", err)
		}

		if r.IsPost() {
			if err == nil {
				port, err := strconv.ParseUint(r.Original().FormValue("port"), 10, 64)

				if err == nil {
					u := bytes.NewBuffer(nil)
					err = xml.EscapeText(u, []byte(r.Original().FormValue("url")))

					if err == nil {
						params := event.NewSetNotificationHTTPHostParamsWithContext(ctx).
							WithHTTPHost(notification.Payload.ID).
							WithHTTPHostNotification(&static.HTTPHostNotification{
								ID:                       notification.Payload.ID,
								ProtocolType:             notification.Payload.ProtocolType,
								ParameterFormatType:      notification.Payload.ParameterFormatType,
								HTTPAuthenticationMethod: notification.Payload.HTTPAuthenticationMethod,

								AddressingFormatType: r.Original().FormValue("address-format"),
								URL:                  &[]string{u.String()}[0],
								PortNo:               port,
							})

						if hostname := r.Original().FormValue("hostname"); len(hostname) > 0 {
							params.HTTPHostNotification.HostName = &hostname
						}

						if ip := r.Original().FormValue("ip"); len(ip) > 0 {
							params.HTTPHostNotification.IPAddress = &ip
						}

						_, err = b.client.Event.SetNotificationHTTPHost(params, nil)
					}
				}

				if err != nil {
					widget.FlashError(r, err, "")
				} else {
					widget.FlashSuccess(r, "Success", "")
				}

				widget.Redirect(r.URL().Path+"?action="+action, http.StatusFound, w, r)
			}

			return
		}

		info := map[string]interface{}{
			"address_format": "",
			"port":           "",
			"url":            "",
			"hostname":       "",
			"ip_address":     "",
		}

		if err == nil {
			info["address_format"] = notification.Payload.AddressingFormatType
			info["port"] = notification.Payload.PortNo

			if notification.Payload.URL != nil {
				info["url"] = *notification.Payload.URL
			}

			if notification.Payload.HostName != nil {
				info["hostname"] = *notification.Payload.HostName
			}

			if notification.Payload.IPAddress != nil {
				info["ip_address"] = *notification.Payload.IPAddress
			}
		}

		vars["notification"] = info
		vars["access_key"] = ""

		if keysConfig := r.Config().String(boggart.ConfigAccessKeys); keysConfig != "" {
			if keys := strings.Split(keysConfig, ","); len(keys) > 0 {
				vars["access_key"] = keys[0]
			}
		}

	case "event":
		if !b.config.EventsEnabled || !r.IsPost() {
			widget.NotFound(w, r)
			return
		}

		if content, err := ioutil.ReadAll(r.Original().Body); err == nil {
			b.Logger().Debug("Call hikvision event " + string(content))

			e := &models.EventNotificationAlert{}

			d := xml.NewDecoder(bytes.NewReader(content))
			d.Strict = false // иногда приходит треш, например кривый амперсанты

			err = d.Decode(e)

			if err != nil {
				b.Logger().Error("Parse event failed",
					"error", err.Error(),
					"body", string(content),
				)

				widget.InternalError(w, r, err)
			} else {
				b.registerEvent(e)
			}
		}

		return
	}

	widget.Render(ctx, "widget", vars)
}

func (b *Bind) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
