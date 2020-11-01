package mqtt

import (
	"context"
	"net"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	subscribers := []mqtt.Subscriber{
		mqtt.NewSubscriber(b.config.TopicDiscoveryPrefix+"/+/"+b.config.TopicPrefix+"/+/config", 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			var component Component

			parts := message.Topic().Split()
			t := ComponentType(parts[len(parts)-4])
			id := parts[len(parts)-2]

			switch t {
			case ComponentTypeBinarySensor:
				component = NewComponentBinarySensor(id)
				// case components.ComponentTypeCover.String():
				// 	component.Type = components.ComponentTypeCover
				// case components.ComponentTypeFan.String():
				// 	component.Type = components.ComponentTypeFan
			case ComponentTypeLight:
				component = NewComponentLight(id)
			case ComponentTypeSensor:
				component = NewComponentSensor(id)
			case ComponentTypeSwitch:
				component = NewComponentSwitch(id)
				// case components.ComponentTypeTextSensor.String():
				// 	component.Type = components.ComponentTypeTextSensor
				// case components.ComponentTypeCamera.String():
				// 	component.Type = components.ComponentTypeCamera
				// case components.ComponentTypeClimate.String():
				// 	component.Type = components.ComponentTypeClimate
			default:
				component = NewComponentBase(id, t)
			}

			if err := message.JSONUnmarshal(component); err != nil {
				return err
			}

			if b.config.IPAddressSensorID == id && b.ipSubscriber.IsFalse() {
				b.MQTT().Subscribe(mqtt.NewSubscriber(component.TopicState(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
					b.ip.Store(net.ParseIP(message.String()))
					return nil
				}))

				b.ipSubscriber.True()
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

	return subscribers
}
