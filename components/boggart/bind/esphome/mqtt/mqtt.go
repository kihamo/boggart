package mqtt

import (
	"context"
	"net"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()

	topicBirth := cfg.TopicBirth
	if topicBirth == "" {
		topicBirth = cfg.TopicPrefix + "/status"
	}

	topicWill := cfg.TopicWill
	if topicWill == "" {
		topicWill = cfg.TopicPrefix + "/status"
	}

	subscribers := []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicDiscoveryPrefix+"/+/"+cfg.TopicPrefix+"/+/config", 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			// компонент удаляют, поэтому пришла пустота на которую реагировать не надо
			if len(message.Payload()) == 0 {
				return nil
			}

			var component Component

			parts := message.Topic().Split()
			t := ComponentType(parts[len(parts)-4])
			id := parts[len(parts)-2]

			switch t {
			case ComponentTypeBinarySensor:
				component = NewComponentBinarySensor(id, message)
				// case components.ComponentTypeCover.String():
				// 	component.Type = components.ComponentTypeCover
				// case components.ComponentTypeFan.String():
				// 	component.Type = components.ComponentTypeFan
			case ComponentTypeLight:
				component = NewComponentLight(id, message)
			case ComponentTypeSensor:
				// text_sensor прикидывается обычным sensor и их надо распознать (HA просто не поддерживает такой тип)
				// распознование очень хрупкое и основывается на проверке признаков unit_of_measurement, expire_after и
				// force_update, которые есть у sensor, но нет у text_sensor
				var check ComponentSensorData
				if err := message.JSONUnmarshal(&check); err == nil {
					if check.UnitOfMeasurement != nil || check.ExpireAfter != nil || check.ForceUpdate != nil {
						component = NewComponentSensor(id, message)
					}
				}

				if component == nil {
					component = NewComponentDefault(id, ComponentTypeTextSensor, message)
				}

			case ComponentTypeSwitch:
				component = NewComponentSwitch(id, message)
				// case components.ComponentTypeTextSensor.String():
				// 	component.Type = components.ComponentTypeTextSensor
				// case components.ComponentTypeCamera.String():
				// 	component.Type = components.ComponentTypeCamera
				// case components.ComponentTypeClimate.String():
				// 	component.Type = components.ComponentTypeClimate
			default:
				component = NewComponentDefault(id, t, message)
			}

			if err := message.JSONUnmarshal(component); err != nil {
				return err
			}

			if t := component.StateTopic(); t != "" {
				component.Subscribe(mqtt.NewSubscriber(component.StateTopic(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					if cfg.IPAddressSensorID == id {
						b.ip.Store(net.ParseIP(message.String()))
					}

					if err := component.SetState(message); err != nil {
						return err
					}

					if cmp, ok := component.(*ComponentBinarySensor); ok && cmp.DeviceClass() == DeviceClassConnectivity && cmp.StateTopic() != "" {
						b.status.Set(cmp.State().(bool))
					}

					return nil
				}))
			}

			return b.register(component)
		}),
		mqtt.NewSubscriber(topicBirth, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if message.String() == cfg.BirthMessage {
				b.status.True()
			}

			return nil
		}),
		mqtt.NewSubscriber(topicWill, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if message.String() == cfg.WillMessage {
				b.status.False()
			}

			return nil
		}),
	}

	return subscribers
}
