package wol

import (
	"context"
	"errors"
	"net"

	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTSubscribers() []mqtt.Subscriber {
	cfg := b.Config().Bind().(*Config)

	return []mqtt.Subscriber{
		mqtt.NewSubscriber(cfg.TopicWOL, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			route := message.Topic().Split()
			if len(route) < 1 {
				return errors.New("bad topic name")
			}

			mac, err := net.ParseMAC(route[len(route)-1])
			if err != nil {
				return err
			}

			return b.WOL(mac, nil, nil)
		}),
		mqtt.NewSubscriber(cfg.TopicWOLWithIPAndSubnet, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			route := message.Topic().Split()
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
