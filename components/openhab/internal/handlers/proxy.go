package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/openhab"
	"github.com/kihamo/boggart/components/openhab/client/things"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/messengers"
)

type ProxyHandler struct {
	dashboard.Handler

	Component  openhab.Component
	Messengers messengers.Component
}

func (h *ProxyHandler) ServeHTTP(w *dashboard.Response, r *dashboard.Request) {
	query := r.URL().Query()

	switch query.Get(":type") {
	case "thing":
		id := query.Get(":id")
		configKey := query.Get(":key")
		if id == "" || configKey == "" {
			h.NotFound(w, r)
			return
		}

		client := h.Component.Client()
		if client == nil {
			http.Error(w, "Client not initialization", http.StatusServiceUnavailable)
			return
		}

		// TODO: cache
		status, err := client.Things.GetByUID(&things.GetByUIDParams{
			ThingUID: id,
			Context:  r.Context(),
		})

		if err != nil {
			break
		}

		configVal, ok := status.Payload.Configuration[configKey]
		if !ok {
			break
		}

		response, err := http.Get(configVal.(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		defer response.Body.Close()

		for headerKey, headerValues := range response.Header {
			for _, headerValue := range headerValues {
				w.Header().Add(headerKey, headerValue)
			}
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		if send := query.Get("send"); r.IsPost() && h.Messengers != nil {
			h.sendFileToMessenger(send, body, w, r)
			return
		}

		w.WriteHeader(response.StatusCode)
		w.Write(body)
		return
	}

	h.NotFound(w, r)
}

func (h *ProxyHandler) sendFileToMessenger(messenger string, body []byte, w *dashboard.Response, r *dashboard.Request) {
	m := h.Messengers.Messenger(messenger)
	if m == nil {
		h.NotFound(w, r)
		return
	}

	var chats []string
	chatsFromConfig := strings.FieldsFunc(r.Config().String(openhab.ConfigTelegramChats), func(c rune) bool {
		return c == ','
	})

	for _, id := range chatsFromConfig {
		chats = append(chats, id)
	}

	if len(chats) == 0 {
		h.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	go func(query url.Values) {
		description := fmt.Sprintf("%s %s at %s", query.Get(":id"), query.Get(":key"), time.Now().Format(time.RFC1123Z))

		switch http.DetectContentType(body[:256]) {
		case "image/jpeg", "image/png", "image/gif", "image/webp":
			for _, chatId := range chats {
				err := m.SendPhoto(chatId, description, bytes.NewReader(body))

				if err != nil {
					return
				}
			}
		default:
			// TODO: send file
		}
	}(r.URL().Query())
}
