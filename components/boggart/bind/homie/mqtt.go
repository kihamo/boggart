package homie

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = "+/+/"

	MQTTPublishTopicBroadcast mqtt.Topic = "+/$broadcast/+"
	MQTTPublishTopicConfigSet            = MQTTPrefix + "$implementation/config/set"
	MQTTPublishTopicReset                = MQTTPrefix + "$implementation/reset"
	MQTTPublishTopicRestart              = MQTTPrefix + "$implementation/restart"

	MQTTSubscribeTopicDeviceAttribute               = MQTTPrefix + "+"
	MQTTSubscribeTopicDeviceAttributeFirmware       = MQTTPrefix + "$fw/+"
	MQTTSubscribeTopicDeviceAttributeImplementation = MQTTPrefix + "$implementation/+/#"
	MQTTSubscribeTopicDeviceAttributeStats          = MQTTPrefix + "$stats/+"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	base := b.config.BaseTopic
	sn := b.SerialNumber()

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicBroadcast.Format(base)),
		mqtt.Topic(MQTTPublishTopicConfigSet.Format(base, sn)),
		mqtt.Topic(MQTTPublishTopicReset.Format(base, sn)),
		mqtt.Topic(MQTTPublishTopicRestart.Format(base, sn)),
	}
}

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	base := b.config.BaseTopic
	sn := b.SerialNumber()

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicDeviceAttribute.Format(base, sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			attributeName := route[len(route)-1]

			if !strings.HasPrefix(attributeName, "$") {
				return errors.New("device attribute name " + attributeName + " is wrong")
			}

			attributeName = attributeName[1:]
			b.registerDeviceAttributes(attributeName, message.String())

			if attributeName == "online" {
				if message.IsTrue() {
					b.UpdateStatus(boggart.BindStatusOnline)
				} else {
					b.UpdateStatus(boggart.BindStatusOffline)
				}
			}

			return nil
		}),
		mqtt.NewSubscriber(MQTTSubscribeTopicDeviceAttributeFirmware.Format(base, sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			b.registerDeviceAttributes("fw."+route[len(route)-1], message.String())
			return nil
		}),
		mqtt.NewSubscriber(MQTTSubscribeTopicDeviceAttributeImplementation.Format(base, sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			// skip special topics
			if message.Topic() == MQTTPublishTopicConfigSet.Format(base, sn) {
				return nil
			}

			route := mqtt.RouteSplit(message.Topic())
			b.registerDeviceAttributes("implementation."+strings.Join(route[3:], "."), message.String())
			return nil
		}),
		mqtt.NewSubscriber(MQTTSubscribeTopicDeviceAttributeStats.Format(base, sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			attributeName := strings.Join(route[3:], ".")
			var value interface{}

			switch attributeName {
			case "interval", "uptime":
				v, _ := strconv.ParseInt(message.String(), 10, 64)
				value = time.Second * time.Duration(v)

			default:
				value = message.String()
			}

			b.registerDeviceAttributes("stats."+attributeName, value)
			return nil
		}),
	}
}
