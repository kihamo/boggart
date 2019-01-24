package grafana

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTSubscribeTopicAnnotation mqtt.Topic = boggart.ComponentName + "/grafana/+/annotation"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	name := mqtt.NameReplace(b.name)

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicAnnotation.Format(name), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			var request struct {
				Title string   `json:"title,omitempty"`
				Text  string   `json:"text,omitempty"`
				Tags  []string `json:"tags,omitempty"`
			}

			if err := message.UnmarshalJSON(&request); err != nil {
				return err
			}

			return b.CreateAnnotation(request.Title, request.Text, request.Tags, nil, nil)
		}),
	}
}
