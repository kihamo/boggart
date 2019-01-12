package gpio

import (
	"bytes"
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go/log"
)

const (
	TopicPinState mqtt.Topic = boggart.ComponentName + "/gpio/+"
	TopicPinSet   mqtt.Topic = boggart.ComponentName + "/gpio/+/set"
)

func (d *Bind) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(TopicPinState.Format(d.pin.Number())),
		mqtt.Topic(TopicPinSet.Format(d.pin.Number())),
	}
}

func (d *Bind) MQTTSubscribers() []mqtt.Subscriber {
	if d.Mode() != GPIOModeOut {
		return nil
	}

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(
			TopicPinSet.Format(d.pin.Number()),
			0,
			func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
				span, ctx := tracing.StartSpanFromContext(ctx, "gpio", "set")
				span.LogFields(
					log.String("name", d.pin.Name()),
					log.Int("number", d.pin.Number()))
				defer span.Finish()

				var err error

				if bytes.Equal(message.Payload(), []byte(`1`)) {
					err = d.High(ctx)
					span.LogFields(log.String("out", "high"))
				} else {
					err = d.Low(ctx)
					span.LogFields(log.String("out", "low"))
				}

				if err != nil {
					tracing.SpanError(span, err)
				}
			}),
	}
}
