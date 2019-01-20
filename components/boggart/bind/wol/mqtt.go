package wol

import (
	"context"
	"errors"
	"net"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTSubscribeTopicWOL                mqtt.Topic = boggart.ComponentName + "/wol/+"
	MQTTSubscribeTopicWOLWithIPAndSubnet mqtt.Topic = boggart.ComponentName + "/wol/+/+/+"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(MQTTSubscribeTopicWOL.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 1 {
				return errors.New("bad topic name")
			}

			mac, err := net.ParseMAC(route[len(route)-1])
			if err != nil {
				return err
			}

			return b.WOL(mac, nil, nil)
		}),
		mqtt.NewSubscriber(MQTTSubscribeTopicWOLWithIPAndSubnet.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			route := mqtt.RouteSplit(message.Topic())
			if len(route) < 3 {
				return errors.New("bad topic name")
			}

			mac, err := net.ParseMAC(route[len(route)-3])
			if err != nil {
				return err
			}

			subnet := net.ParseIP(route[len(route)-1])
			ip := net.ParseIP(route[len(route)-2])

			return b.WOL(mac, ip, subnet)
		}),
	}
}
