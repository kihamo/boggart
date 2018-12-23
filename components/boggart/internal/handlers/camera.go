package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/messengers"
)

const (
	queryParamSendToMessenger = "send"
)

type CameraHandler struct {
	dashboard.Handler

	DevicesManager boggart.DevicesManager
	Messengers     messengers.Component
}

func (h *CameraHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	sn := q.Get(":sn")
	channel := q.Get(":channel")

	if sn == "" || channel == "" {
		h.NotFound(w, r)
		return
	}

	ch, err := strconv.ParseUint(channel, 10, 64)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	sn = strings.Replace(sn, "/", "_", -1)

	devicesList := h.DevicesManager.DevicesByTypes([]boggart.DeviceType{boggart.DeviceTypeCamera})
	if len(devicesList) == 0 {
		h.NotFound(w, r)
		return
	}

	var device *devices.CameraHikVision

	for _, d := range devicesList {
		hv, ok := d.(*devices.CameraHikVision)
		if !ok {
			continue
		}

		if strings.Replace(hv.SerialNumber(), "/", "_", -1) != sn {
			continue
		}

		device = hv
		break
	}

	if device == nil {
		h.NotFound(w, r)
		return
	}

	if messenger := q.Get(queryParamSendToMessenger); r.IsPost() && h.Messengers != nil {
		m := h.Messengers.Messenger(messenger)
		if m == nil {
			h.NotFound(w, r)
			return
		}

		var chats []string
		chatsFromConfig := strings.FieldsFunc(r.Config().String(boggart.ConfigListenerTelegramChats), func(c rune) bool {
			return c == ','
		})

		for _, id := range chatsFromConfig {
			chats = append(chats, id)
		}

		if len(chats) == 0 {
			h.NotFound(w, r)
			return
		}

		buf := &bytes.Buffer{}
		if err := device.Snapshot(r.Context(), ch, buf); err != nil {
			h.NotFound(w, r)
		}

		w.WriteHeader(http.StatusNoContent)

		go func(query url.Values) {
			description := fmt.Sprintf("%s at %s", device.Description(), time.Now().Format(time.RFC1123Z))
			mime := http.DetectContentType(buf.Bytes()[:256])

			switch mime {
			case "image/jpeg", "image/png", "image/gif", "image/webp":
				for _, chatId := range chats {
					if err := m.SendPhoto(chatId, description, buf); err != nil {
						h.Logger().Error("Failed send message to "+messenger, "error", err.Error())
						return
					}
				}
			default:
				// TODO: send file
				h.Logger().Warn("Send file not implemented", "mime", mime)
			}
		}(q)

		return
	}

	err = device.Snapshot(r.Context(), ch, w)
	if err != nil {
		h.NotFound(w, r)
	}
}
