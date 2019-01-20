package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/hikvision"
	"github.com/kihamo/boggart/components/boggart/internal/manager"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/messengers"
)

const (
	queryParamSendToMessenger = "send"
)

type CameraHandler struct {
	dashboard.Handler

	DevicesManager *manager.Manager
	Messengers     messengers.Component
}

func (h *CameraHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	q := r.URL().Query()

	id := q.Get(":id")
	channel := q.Get(":channel")

	if id == "" || channel == "" {
		h.NotFound(w, r)
		return
	}

	ch, err := strconv.ParseUint(channel, 10, 64)
	if err != nil {
		h.NotFound(w, r)
		return
	}

	bindItem := h.DevicesManager.Bind(id)
	if bindItem == nil {
		h.NotFound(w, r)
		return
	}

	bind, ok := bindItem.Bind().(*hikvision.Bind)
	if !ok {
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
		if err := bind.Snapshot(r.Context(), ch, buf); err != nil {
			h.Logger().Warn("Failed get snapshot", "error", err, "device", bindItem.Description())

			h.NotFound(w, r)
			return
		}

		w.WriteHeader(http.StatusNoContent)

		go func(query url.Values) {
			description := fmt.Sprintf("%s at %s", bindItem.Description(), time.Now().Format(time.RFC1123Z))

			var mime string

			if buf.Len() < 256 {
				mime = http.DetectContentType(buf.Bytes())
			} else {
				mime = http.DetectContentType(buf.Bytes()[:256])
			}

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

	buf := bytes.NewBuffer(nil)
	if err = bind.Snapshot(r.Context(), ch, buf); err != nil {
		h.NotFound(w, r)
		return
	}

	io.Copy(w, buf)
}
