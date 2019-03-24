package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	httpClient "github.com/kihamo/boggart/components/boggart/protocols/http"
	"github.com/kihamo/boggart/components/openhab"
	"github.com/kihamo/boggart/components/openhab/client/things"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/kihamo/shadow/components/tracing"
)

const (
	queryParamSendToMessenger = "send"
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

		span, ctx := tracing.StartSpanFromContext(r.Context(), h.Component.Name(), "api.things.GetByUID")

		// TODO: cache
		status, err := client.Things.GetByUID(&things.GetByUIDParams{
			ThingUID: id,
			Context:  ctx,
		})

		if err != nil {
			tracing.SpanError(span, err)
			break
		}

		configVal, ok := status.Payload.Configuration[configKey]
		if !ok {
			tracing.SpanError(span, fmt.Errorf("%s not found in response", configKey))
			break
		}

		span.Finish()

		h.proxy(configVal.(string), w, r)
		return
	}

	h.NotFound(w, r)
}

func (h *ProxyHandler) proxy(u string, w *dashboard.Response, r *dashboard.Request) {
	span, ctx := tracing.StartSpanFromContext(r.Context(), h.Component.Name(), "handler.proxy")
	defer span.Finish()

	response, err := httpClient.NewClient().Get(ctx, u)
	if err != nil {
		tracing.SpanError(span, err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tracing.SpanError(span, err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	if send := r.URL().Query().Get(queryParamSendToMessenger); r.IsPost() && h.Messengers != nil {
		h.sendFileToMessenger(send, body, w, r)
		return
	}

	for headerKey, headerValues := range response.Header {
		for _, headerValue := range headerValues {
			w.Header().Add(headerKey, headerValue)
		}
	}

	w.WriteHeader(response.StatusCode)
	_, _ = w.Write(body)
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

	chats = append(chats, chatsFromConfig...)

	if len(chats) == 0 {
		h.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	go func(query url.Values) {
		description := fmt.Sprintf("%s %s at %s", query.Get(":id"), query.Get(":key"), time.Now().Format(time.RFC1123Z))
		mime := http.DetectContentType(body[:256])

		switch mime {
		case "image/jpeg", "image/png", "image/gif", "image/webp":
			for _, chatId := range chats {
				if err := m.SendPhoto(chatId, description, bytes.NewReader(body)); err != nil {
					logging.Log(r.Context()).Error("Failed send message to "+messenger, "error", err.Error())
					return
				}
			}
		default:
			// TODO: send file
			logging.Log(r.Context()).Warn("Send file not implemented", "mime", mime)
		}
	}(r.URL().Query())
}
