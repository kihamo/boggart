package mqtt

import (
	"context"
	"net"
	"net/url"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.config()

	topicLog := cfg.TopicLog
	if topicLog == "" {
		topicLog = cfg.TopicPrefix + "/debug"
	}

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
			var component Component

			parts := message.Topic().Split()
			t := ComponentType(parts[len(parts)-4])
			id := parts[len(parts)-2]

			// компонент удаляют, поэтому пришла пустота на которую реагировать не надо
			if len(message.Payload()) == 0 {
				b.delete(id)
				return nil
			}

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
					if len(message.Payload()) == 0 {
						return nil
					}

					if cfg.IPAddressSensorID == id {
						ip := net.ParseIP(message.String())
						b.ip.Store(ip)

						b.Meta().SetLink(&url.URL{
							Scheme: "http",
							Host:   ip.String(),
						})
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
		//mqtt.NewSubscriber(topicLog, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
		//	/*
		//		   "",    // NONE
		//		   "E",   // ERROR
		//		   "W",   // WARNING
		//		   "I",   // INFO
		//		   "C",   // CONFIG
		//		   "D",   // DEBUG
		//		   "V",   // VERBOSE
		//		   "VV",  // VERY_VERBOSE
		//
		//			this->printf_to_buffer_("%s[%s][%s:%03u]: ", color, letter, tag, line);
		//
		//		[0;36m[D][sensor:113]: 'Temperature DS': Sending state 21.62500 °C with 1 decimals of accuracy[0m
		//	*/
		//
		//	fmt.Println(message)
		//	fmt.Println(message.Payload())
		//	fmt.Println(string([]byte{27, 91, 48, 59, 51, 54}))
		//	fmt.Println(hex.EncodeToString([]byte{27, 91, 48, 59, 51, 54}))
		//	fmt.Println(message.HEX())
		//	fmt.Println()
		//
		//	return nil
		//}),
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
