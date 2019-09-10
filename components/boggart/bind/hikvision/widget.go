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
	cfg := b.Config().(*Config)

	vars := map[string]interface{}{
		"action":                   action,
		"events_enabled":           cfg.EventsEnabled,
		"preview_refresh_interval": cfg.PreviewRefreshInterval.Seconds(),
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

	case "notification":
		if !bind.config.EventsEnabled {
			t.NotFound(w, r)
			return
		}

		params := event.NewGetNotificationHttpHostParamsWithContext(ctx).
			WithHttpHost(1)

		notification, err := bind.client.Event.GetNotificationHttpHost(params, nil)
		if err != nil {
			r.Session().FlashBag().Error(t.Translate(ctx, "Get notification http host failed with error %s", "", err.Error()))
		}

		if r.IsPost() {
			if err == nil {
				port, err := strconv.ParseUint(r.Original().FormValue("port"), 10, 64)

				u := bytes.NewBuffer(nil)
				err = xml.EscapeText(u, []byte(r.Original().FormValue("url")))

				if err == nil {
					params := event.NewSetNotificationHttpHostParamsWithContext(ctx).
						WithHttpHost(notification.Payload.ID).
						WithHttpHostNotification(&models.HttpHostNotification{
							ID:                       notification.Payload.ID,
							ProtocolType:             notification.Payload.ProtocolType,
							ParameterFormatType:      notification.Payload.ParameterFormatType,
							HttpAuthenticationMethod: notification.Payload.HttpAuthenticationMethod,

							AddressingFormatType: r.Original().FormValue("address-format"),
							URL:                  &[]string{u.String()}[0],
							PortNo:               port,
						})

					if hostname := r.Original().FormValue("hostname"); len(hostname) > 0 {
						params.HttpHostNotification.HostName = &hostname
					}

					if ip := r.Original().FormValue("ip"); len(ip) > 0 {
						params.HttpHostNotification.IPAddress = &ip
					}

					_, err = bind.client.Event.SetNotificationHttpHost(params, nil)
				}

				if err != nil {
					r.Session().FlashBag().Error(err.Error())
				} else {
					r.Session().FlashBag().Success(t.Translate(ctx, "Success", ""))
				}

				t.Redirect(r.URL().Path+"?action="+action, http.StatusFound, w, r)
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
		if !bind.config.EventsEnabled || !r.IsPost() {
			t.NotFound(w, r)
			return
		}

		if content, err := ioutil.ReadAll(r.Original().Body); err == nil {
			bind.Logger().Debug("Call hikvision event " + string(content))

			e := &models.EventNotificationAlert{}

			d := xml.NewDecoder(bytes.NewReader(content))
			d.Strict = false // иногда приходит треш, например кривый амперсанты

			err = d.Decode(e)

			if err != nil {
				bind.Logger().Error("Parse event failed",
					"error", err.Error(),
					"body", string(content),
				)

				t.InternalError(w, r, err)
			} else {
				bind.registerEvent(e)
			}
		}

		return
	}

	t.Render(ctx, "widget", vars)
}

func (t Type) WidgetAssetFS() *assetfs.AssetFS {
	return assetFS()
}
