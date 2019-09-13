package grafana

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicAnnotation, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
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
