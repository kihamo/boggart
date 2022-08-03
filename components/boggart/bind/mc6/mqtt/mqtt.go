package mqtt

import (
	"context"
	"fmt"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()
	id := b.Meta().ID()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicMC6Update, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			update, err := b.ParseUpdate(message.Payload())
			if err != nil {
				return err
			}

			if update.MAC != nil && b.Meta().MAC() == nil {
				v := *update.MAC
				err = b.Meta().SetMACAsString(fmt.Sprintf("%s:%s:%s:%s:%s:%s", v[0:2], v[2:4], v[4:6], v[6:8], v[8:10], v[10:12]))
				if err != nil {
					return err
				}
			}

			if update.Temperature != nil {
				metricTemperature.With("id", id).Set(float64(*update.Temperature) / 10)
			}

			if update.HoldTemperature != nil {
				metricHoldTemperature.With("id", id).Set(float64(*update.HoldTemperature) / 10)
			}

			if update.Humidity != nil {
				metricHumidity.With("id", id).Set(float64(*update.Humidity) / 10)
			}

			return err
		}),
		mqtt.NewSubscriber(cfg.TopicSetTemperature, 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			mac := b.Meta().MAC()
			if mac == nil {
				return nil
			}

			//update := &Update{
			//	SetTemperature: &[]int64{message.Int64()}[0],
			//}

			// return b.MQTT().Publish(ctx, cfg.TopicMC6SetTemperature.Format(mac.String(), ))
			return nil
		}),
	}
}
