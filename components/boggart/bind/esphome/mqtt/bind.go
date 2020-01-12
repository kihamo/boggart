package mqtt

import (
	"errors"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/esphome/mqtt/components"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config     *Config
	components sync.Map
	status     atomic.BoolNull
}

func (b *Bind) Close() (err error) {
	if client := b.MQTTClient(); client != nil {
		for _, component := range b.Components() {
			if err = client.UnsubscribeSubscribers(component.Subscribers()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Bind) Component(id string) (cmp components.Component) {
	b.components.Range(func(_, value interface{}) bool {
		if c := value.(components.Component); c.GetID() == id {
			cmp = c
			return false
		}

		return true
	})

	return cmp
}

func (b *Bind) Components() []components.Component {
	list := make([]components.Component, 0)

	b.components.Range(func(_, value interface{}) bool {
		list = append(list, value.(components.Component))
		return true
	})

	return list
}

func (b *Bind) register(c components.Component) (err error) {
	b.components.Store(c.GetUniqueID(), c)

	subscribers := c.Subscribers()
	if len(subscribers) > 0 {
		client := b.MQTTClient()

		if client == nil {
			return errors.New("MQTT client isn't init")
		}

		for _, subscribe := range subscribers {
			if err = client.SubscribeSubscriber(subscribe); err != nil {
				return err
			}
		}
	}

	return err
}
