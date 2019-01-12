package broadlink

import (
	"bytes"
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go/log"
)

const (
	SP3SUpdateInterval = time.Second * 3 // as e-control app, refresh every 3 sec

	SP3SMQTTTopicState mqtt.Topic = boggart.ComponentName + "/socket/+/state"
	SP3SMQTTTopicPower mqtt.Topic = boggart.ComponentName + "/socket/+/power"
	SP3SMQTTTopicSet   mqtt.Topic = boggart.ComponentName + "/socket/+/set"
)

func (b *BindSP3S) MQTTTopics() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Topic{
		mqtt.Topic(SP3SMQTTTopicState.Format(sn)),
		mqtt.Topic(SP3SMQTTTopicPower.Format(sn)),
		mqtt.Topic(SP3SMQTTTopicSet.Format(sn)),
	}
}

func (b *BindSP3S) MQTTSubscribers() []mqtt.Subscriber {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(SP3SMQTTTopicSet.Format(sn), 0, func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
			if b.Status() != boggart.DeviceStatusOnline {
				return
			}

			span, ctx := tracing.StartSpanFromContext(ctx, "socket", "set")
			span.LogFields(
				log.String("mac", b.provider.MAC().String()),
				log.String("ip", b.provider.Addr().String()))
			defer span.Finish()

			var err error

			if bytes.Equal(message.Payload(), []byte(`1`)) {
				err = b.On(ctx)
				span.LogFields(log.String("state", "on"))
			} else {
				err = b.Off(ctx)
				span.LogFields(log.String("state", "off"))
			}

			if err != nil {
				tracing.SpanError(span, err)
			}
		}),
	}
}
