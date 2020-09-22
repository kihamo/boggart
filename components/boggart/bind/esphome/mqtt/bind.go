package mqtt

import (
	"net"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
)

type Bind struct {
	di.MetaBind
	di.LoggerBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	config     *Config
	components sync.Map
	status     atomic.BoolNull

	ip      net.IP
	ipMutex sync.RWMutex
}

func (b *Bind) Close() (err error) {
	client := b.MQTT()

	for _, component := range b.Components() {
		for _, subscribe := range component.Subscribers() {
			if err = client.Unsubscribe(subscribe); err != nil {
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

func (b *Bind) register(component Component) (err error) {
	b.components.Store(component.GetUniqueID(), component)

	if mac := component.GetDevice().MAC(); mac != nil {
		b.Meta().SetMAC(mac)
	}

	subscribers := component.Subscribers()
	if len(subscribers) > 0 {
		client := b.MQTT()

		for _, subscribe := range subscribers {
			if err = client.Subscribe(subscribe); err != nil {
				return err
			}
		}
	}

	return err
}
