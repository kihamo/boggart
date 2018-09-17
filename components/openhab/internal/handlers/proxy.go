package handlers

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/openhab"
	"github.com/kihamo/boggart/components/openhab/client/things"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/messengers"
	"github.com/mattn/go-mjpeg"
)

const (
	queryParamSendToMessenger = "send"
	queryParamStream          = "stream"

	streamMJPEG = "mjpeg"
)

var (
	mjpegLock    sync.RWMutex
	mjpegStreams = map[string]*mjpeg.Stream{}
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

		switch strings.ToLower(query.Get(queryParamStream)) {
		case streamMJPEG:
			h.proxymjpegStream(configVal.(string), w, r)

		default:
			h.proxy(configVal.(string), w, r)
		}

		return
	}

	h.NotFound(w, r)
}

func (h *ProxyHandler) proxymjpegStream(u string, w *dashboard.Response, r *dashboard.Request) {
	var stream *mjpeg.Stream

	mjpegLock.RLock()
	stream, ok := mjpegStreams[u]
	mjpegLock.RUnlock()

	if !ok {
		dec, err := mjpeg.NewDecoderFromURL(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		stream = mjpeg.NewStreamWithInterval(r.Config().Duration(openhab.ConfigProxyMJPEGInterval))

		go func(u string, l logger.Logger) {
			defer func() {
				mjpegLock.Lock()
				delete(mjpegStreams, u)
				mjpegLock.Unlock()

				l.Warnf("MJPEG stream for %s is down", u)
			}()

			var buf bytes.Buffer
			for {
				img, err := dec.Decode()
				if err != nil {
					l.Errorf("MJPEG stream for %s MJPEG decode failed", u, map[string]interface{}{
						"error": err.Error(),
					})
					break
				}
				buf.Reset()
				err = jpeg.Encode(&buf, img, nil)
				if err != nil {
					l.Errorf("MJPEG stream for %s JPEG encode failed", u, map[string]interface{}{
						"error": err.Error(),
					})
					break
				}
				err = stream.Update(buf.Bytes())
				if err != nil {
					l.Errorf("MJPEG stream for %s stream update failed", u, map[string]interface{}{
						"error": err.Error(),
					})
					break
				}
			}
		}(u, dashboard.LoggerFromContext(r.Context()))

		mjpegLock.Lock()
		mjpegStreams[u] = stream
		mjpegLock.Unlock()
	}

	if send := r.URL().Query().Get(queryParamSendToMessenger); r.IsPost() && h.Messengers != nil {
		h.sendFileToMessenger(send, stream.Current(), w, r)
		return
	}

	stream.ServeHTTP(w, r.Original())
}

func (h *ProxyHandler) proxy(u string, w *dashboard.Response, r *dashboard.Request) {
	response, err := http.Get(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
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
	w.Write(body)
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
		mime := http.DetectContentType(body[:256])

		switch mime {
		case "image/jpeg", "image/png", "image/gif", "image/webp":
			for _, chatId := range chats {
				if err := m.SendPhoto(chatId, description, bytes.NewReader(body)); err != nil {
					dashboard.LoggerFromContext(r.Context()).Errorf("Failed send message to %s", messenger, map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
			}
		default:
			// TODO: send file
			dashboard.LoggerFromContext(r.Context()).Warn("Send file not implemented", map[string]interface{}{
				"mime": mime,
			})
		}
	}(r.URL().Query())
}
