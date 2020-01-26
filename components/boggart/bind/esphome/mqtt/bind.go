package mqtt

import (
	"errors"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.ProbesBind

	config     *Config
	components sync.Map
	status     atomic.BoolNull
}

func (b *Bind) Close() (err error) {
	if client := b.MQTT().Client(); client != nil {
		for _, component := range b.Components() {
			if err = client.UnsubscribeSubscribers(component.Subscribers()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Bind) Component(id string) (cmp Component) {
	b.components.Range(func(_, value interface{}) bool {
		if c := value.(Component); c.GetID() == id {
			cmp = c
			return false
		}

		return true
	})

	return cmp
}

func (b *Bind) Components() []Component {
	list := make([]Component, 0)

	b.components.Range(func(_, value interface{}) bool {
		list = append(list, value.(Component))
		return true
	})

	return list
}

func (b *Bind) register(c Component) (err error) {
	b.components.Store(c.GetUniqueID(), c)

	if mac := c.GetDevice().MAC(); mac != nil {
		b.Meta().SetMAC(mac)
	}

	subscribers := c.Subscribers()
	if len(subscribers) > 0 {
		client := b.MQTT().Client()

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
