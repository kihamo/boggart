package mqtt

import (
	"context"
	"fmt"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()
	id := b.Meta().ID()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicMC6Update, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			mac := b.Meta().MACAsString()
			if mac != "" {
				if !b.MQTT().CheckValueInTopic(message.Topic(), strings.ReplaceAll(mac, ":", ""), -1) {
					return nil
				}
			}

			update, err := b.ParseUpdate(message.Payload())
			if err != nil {
				return err
			}

			if mac == "" && update.MAC != nil {
				v := *update.MAC
				err = b.Meta().SetMACAsString(fmt.Sprintf("%s:%s:%s:%s:%s:%s", v[0:2], v[2:4], v[4:6], v[6:8], v[8:10], v[10:12]))
				if err != nil {
					return err
				}
			}

			if update.Temperature != nil {
				metricTemperature.With("id", id).Set(float64(*update.Temperature) / 10)
			}

			if update.SetTemperature != nil {
				metricSetTemperature.With("id", id).Set(float64(*update.SetTemperature) / 10)
			}

			if update.Humidity != nil {
				metricHumidity.With("id", id).Set(float64(*update.Humidity) / 10)
			}

			return err
		}),
		mqtt.NewSubscriber(cfg.TopicSetTemperature.Format(id), 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			mac := b.Meta().MACAsString()
			if mac == "" {
				return nil
			}
			mac = strings.ReplaceAll(mac, ":", "")

			payload, err := b.GenerateUpdate(&Update{
				SetTemperature: &[]int64{int64(message.Float64() * 10)}[0],
			})

			if err != nil {
				return err
			}

			return b.MQTT().Publish(ctx, cfg.TopicMC6SetTemperature.Format(mac), payload)
		}),
	}
}
