package modbus

import (
	"context"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()
	id := b.Meta().ID()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicPower.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Power(ctx, message.Bool())
		}),
		mqtt.NewSubscriber(cfg.TopicSetTemperature.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Provider().SetTemperature(message.Float64())
		}),
		mqtt.NewSubscriber(cfg.TopicAway.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.Away(ctx, message.Bool())
		}),
		//mqtt.NewSubscriber(cfg.TopicTemperatureFormat.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
		//	return b.TemperatureFormat(ctx, uint16(message.Uint64()))
		//}),
		mqtt.NewSubscriber(cfg.TopicAwayTemperature.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.AwayTemperature(ctx, uint16(message.Uint64()))
		}),
		mqtt.NewSubscriber(cfg.TopicHoldingTemperature.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.HoldingTemperature(ctx, uint16(message.Uint64()))
		}),
	}
}
