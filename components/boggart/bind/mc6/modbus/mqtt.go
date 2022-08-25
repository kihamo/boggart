package modbus

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()
	id := b.Meta().ID()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicPower.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetStatus(message.Bool())
			}),
		),
		mqtt.NewSubscriber(cfg.TopicTargetTemperature.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetTargetTemperature(message.Float64())
			}),
		),
		mqtt.NewSubscriber(cfg.TopicAway.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetAway(message.Bool())
			}),
		),
		mqtt.NewSubscriber(cfg.TopicAwayTemperature.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetAwayTemperature(message.Float64())
			}),
		),
		mqtt.NewSubscriber(cfg.TopicHoldingTemperature.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetHoldingTemperature(message.Float64())
			}),
		),
	}
}

func (b *Bind) MQTTWrapSubscriberChangeState(f func(message mqtt.Message) error) func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
	return func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
		err := f(message)

		if err == nil {
			return b.notifyChangeState(ctx)
		}

		return err
	}
}
