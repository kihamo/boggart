package grafana

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/grafana/client/annotation"
	"go.uber.org/multierr"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicAnnotation, 0, b.MQTT().WrapSubscribeDeviceIsOnline(func(ctx context.Context, _ mqtt.Component, message mqtt.Message) (err error) {
			var request struct {
				Title string   `json:"title,omitempty"`
				Text  string   `json:"text,omitempty"`
				Tags  []string `json:"tags,omitempty"`
			}

			if err := message.JSONUnmarshal(&request); err != nil {
				return err
			}

			params := annotation.NewCreateAnnotationParamsWithContext(ctx).
				WithRequest(annotation.CreateAnnotationBody{
					Time: time.Now().UnixNano() / int64(time.Millisecond),
					Text: request.Text,
					Tags: request.Tags,
				})

			for _, dashboardID := range b.config.Dashboards {
				params.Request.DashboardID = dashboardID

				if _, e := b.provider.Annotation.CreateAnnotation(params, nil); e != nil {
					err = multierr.Append(err, e)
				}
			}

			return err
		})),
	}
}
