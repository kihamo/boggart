package telegram

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/mime"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicSendMessage, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTSendMessage)),
		mqtt.NewSubscriber(cfg.TopicSendFile, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTSendFile)),
		mqtt.NewSubscriber(cfg.TopicSendFileURL, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTSendFileURL)),
		mqtt.NewSubscriber(cfg.TopicSendFileBase64, 0, b.MQTT().WrapSubscribeDeviceIsOnline(b.callbackMQTTSendFileBase64)),
	}
}

func (b *Bind) callbackMQTTSendMessage(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -4) {
		return nil
	}

	routes := message.Topic().Split()
	if len(routes) < 1 {
		return errors.New("bad topic name")
	}

	return b.SendMessage(routes[len(routes)-2], message.String())
}

func (b *Bind) callbackMQTTSendFile(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -4) {
		return nil
	}

	routes := message.Topic().Split()
	if len(routes) < 1 {
		return errors.New("bad topic name")
	}

	var (
		mimeType mime.Type
		name     string
		url      string

		payload FilePayload
	)

	if err := message.JSONUnmarshal(&payload); err == nil {
		mimeType = mime.Type(payload.MIME)
		name = payload.Name
		url = payload.URL
	} else {
		url = message.String()
	}

	if url == "" {
		return errors.New("url fields is empty")
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if mimeType == mime.TypeUnknown {
		mimeType, err = mime.TypeFromHTTPHeader(response.Header)
		if err != nil {
			var restored io.Reader

			mimeType, restored, err = mime.TypeFromDataRestored(response.Body)
			if err != nil {
				return err
			}

			response.Body = ioutil.NopCloser(restored)
		}
	}

	to := routes[len(routes)-2]

	if name == "" {
		name = "File at " + time.Now().Format(time.RFC1123Z)
	}

	switch v := mimeType; {
	case v.IsImage():
		err = b.SendPhoto(to, name, response.Body)
	case v.IsAudio():
		err = b.SendAudio(to, name, response.Body)
	default:
		err = b.SendDocument(to, name, response.Body)
	}

	return err
}

func (b *Bind) callbackMQTTSendFileURL(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -5) {
		return nil
	}

	routes := message.Topic().Split()
	if len(routes) < 1 {
		return errors.New("bad topic name")
	}

	return b.SendFileAsURL(routes[len(routes)-3], "File at "+time.Now().Format(time.RFC1123Z), message.String())
}

func (b *Bind) callbackMQTTSendFileBase64(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	if !b.MQTT().CheckSerialNumberInTopic(message.Topic(), -5) || message.Len() == 0 {
		return nil
	}

	routes := message.Topic().Split()
	if len(routes) < 1 {
		return errors.New("bad topic name")
	}

	payload, err := message.Base64()
	if err != nil {
		return err
	}

	mimeType, restored, err := mime.TypeFromDataRestored(bytes.NewReader(payload))
	if err != nil {
		return err
	}

	to := routes[len(routes)-3]
	name := "File at " + time.Now().Format(time.RFC1123Z)

	switch v := mimeType; {
	case v.IsImage():
		err = b.SendPhoto(to, name, restored)
	case v.IsAudio():
		err = b.SendAudio(to, name, restored)
	default:
		err = b.SendDocument(to, name, restored)
	}

	return err
}
