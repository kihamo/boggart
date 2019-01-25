package telegram

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/storage"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/telegram/+/"

	MQTTSubscribeTopicMessage = MQTTPrefix + "message/+"
	MQTTSubscribeTopicFile    = MQTTPrefix + "file/+"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicMessage.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 3) {
				return nil
			}

			routes := mqtt.RouteSplit(message.Topic())
			if len(routes) < 1 {
				return errors.New("bad topic name")
			}

			return b.SendMessage(routes[len(routes)-1], message.String())
		}),
		mqtt.NewSubscriber(MQTTSubscribeTopicFile.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if !boggart.CheckSerialNumberInMQTTTopic(b, message.Topic(), 3) {
				return nil
			}

			routes := mqtt.RouteSplit(message.Topic())
			if len(routes) < 1 {
				return errors.New("bad topic name")
			}

			mimeType, err := storage.MimeTypeFromURL(message.String())
			if err != nil {
				return err
			}

			request, err := http.NewRequest(http.MethodGet, message.String(), nil)
			if err != nil {
				return err
			}

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				return err
			}
			defer response.Body.Close()

			to := routes[len(routes)-1]
			name := "File at " + time.Now().Format(time.RFC1123Z)

			switch mimeType {
			case storage.MIMETypeJPEG:
				err = b.SendPhoto(to, name, response.Body)

			case storage.MIMETypeMPEG, storage.MIMETypeWAVE, storage.MIMETypeOGG:
				err = b.SendAudio(to, name, response.Body)

			default:
				err = b.SendDocument(to, name, response.Body)
			}

			return err
		}),
	}
}
