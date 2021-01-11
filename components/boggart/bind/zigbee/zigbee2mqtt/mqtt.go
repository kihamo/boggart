package zigbee2mqtt

import (
	"bytes"
	"context"
	"errors"

	"github.com/kihamo/boggart/components/mqtt"
)

var (
	PayloadOnline  = []byte(`online`)
	PayloadOffline = []byte(`offline`)
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	subscribers := []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicState, 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.setState(ctx, bytes.Equal(message.Payload(), PayloadOnline))
		}),
		mqtt.NewSubscriber(b.config.TopicConfig, 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
			return b.setSettingsFromMessage(message)
		}),
	}

	if b.config.NewAPI {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(b.config.TopicDevices, 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
				var devices []*DeviceNewAPI

				if err := message.JSONUnmarshal(&devices); err != nil {
					return err
				}

				b.devicesLock.Lock()
				b.devices = make(map[string]*Device, len(devices))

				for _, d := range devices {
					b.devices[d.IEEEAddress] = &Device{
						FriendlyName:   d.FriendlyName,
						IEEEAddress:    d.IEEEAddress,
						NetworkAddress: d.NetworkAddress,
						Type:           d.Type,
					}
				}

				b.devicesLock.Unlock()

				return nil
			}),
			mqtt.NewSubscriber(b.config.TopicHealthCheckResponse, 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
				var hc HealthCheck

				if err := message.JSONUnmarshal(&hc); err != nil {
					return err
				}

				return b.setState(ctx, hc.Status == "ok" && hc.Data.Healthy)
			}),
			mqtt.NewSubscriber(b.config.TopicLogging, 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
				var log LogNewAPI

				if err := message.JSONUnmarshal(&log); err != nil {
					return err
				}

				var logger func(string, ...interface{})

				switch log.Level {
				case "debug":
					logger = b.Logger().Debug
				case "info":
					logger = b.Logger().Info
				case "warn":
					logger = b.Logger().Warn
				case "error":
					logger = b.Logger().Error
				default:
					return errors.New("unknown log level " + log.Level)
				}

				logger(log.Message)

				return nil
			}),
			// info отправляется один раз, config при старте + при изменениях, поэтому config надо слушать всегда
			mqtt.NewSubscriber(b.config.TopicInfo, 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
				return b.setSettingsFromMessage(message)
			}),
		)
	} else {
		subscribers = append(subscribers,
			mqtt.NewSubscriber(b.config.TopicDevicesResponse, 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
				var devices []*Device

				if err := message.JSONUnmarshal(&devices); err != nil {
					return err
				}

				b.devicesLock.Lock()
				b.devices = make(map[string]*Device, len(devices))

				for _, d := range devices {
					b.devices[d.IEEEAddress] = d
				}

				b.devicesLock.Unlock()

				return nil
			}),
			mqtt.NewSubscriber(b.config.TopicLog, 0, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
				var log Log

				if err := message.JSONUnmarshal(&log); err != nil {
					return err
				}

				for _, msg := range log.Message {
					if m, ok := msg.(map[string]interface{}); ok {
						fields := make([]interface{}, 0, len(m)*2)

						for k, v := range m {
							fields = append(fields, k, v)
						}

						b.Logger().Info(log.Type, fields...)
					} else {
						b.Logger().Info(log.Type, msg)
					}
				}

				return nil
			}),
		)
	}

	return subscribers
}
