package mqtt

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/bind/esphome/mqtt/components"
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicDiscoveryPrefix+"/+/"+b.config.TopicPrefix+"/+/config", 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			var component components.Component

			parts := message.Topic().Split()
			t := parts[len(parts)-4]
			id := parts[len(parts)-2]

			switch t {
			//case components.ComponentTypeBinarySensor.String():
			//	component.Type = components.ComponentTypeBinarySensor
			//case components.ComponentTypeCover.String():
			//	component.Type = components.ComponentTypeCover
			//case components.ComponentTypeFan.String():
			//	component.Type = components.ComponentTypeFan
			//case components.ComponentTypeLight.String():
			//	component.Type = components.ComponentTypeLight
			case components.ComponentTypeSensor.String():
				component = &components.Sensor{
					Base: &components.Base{
						ID:   id,
						Type: components.ComponentTypeSensor,
					},
				}
			//case components.ComponentTypeSwitch.String():
			//	component.Type = components.ComponentTypeSwitch
			//case components.ComponentTypeTextSensor.String():
			//	component.Type = components.ComponentTypeTextSensor
			//case components.ComponentTypeCamera.String():
			//	component.Type = components.ComponentTypeCamera
			//case components.ComponentTypeClimate.String():
			//	component.Type = components.ComponentTypeClimate
			default:
				component = &components.Base{
					ID:   id,
					Type: components.ComponentType(t),
				}
			}

			if err := message.UnmarshalJSON(component); err != nil {
				return err
			}

			return b.register(component)
		}),
		mqtt.NewSubscriber(b.config.TopicBirth, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if message.String() == b.config.BirthMessage {
				b.status.True()
			}

			return nil
		}),
		mqtt.NewSubscriber(b.config.TopicWill, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			if message.String() == b.config.WillMessage {
				b.status.False()
			}

			return nil
		}),
	}
}
