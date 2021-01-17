package mqtt

import (
	"context"
	"net"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Bind struct {
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	config                 *Config
	components             sync.Map
	status                 atomic.BoolNull
	connectivitySubscriber *atomic.Bool

	ip           *atomic.Value
	ipSubscriber *atomic.Bool
}

func (b *Bind) Run() error {
	b.status.Nil()
	b.ipSubscriber.False()
	b.connectivitySubscriber.False()

	return nil
}

func (b *Bind) Component(id string) (cmp Component) {
	b.components.Range(func(_, value interface{}) bool {
		if c := value.(Component); c.ID() == id {
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

func (b *Bind) IP() net.IP {
	if v := b.ip.Load(); v != nil {
		return v.(net.IP)
	}

	return nil
}

func (b *Bind) register(component Component) (err error) {
	if _, ok := b.components.Load(component.UniqueID()); ok {
		return nil
	}

	b.components.Store(component.UniqueID(), component)

	if mac := component.Device().MAC(); mac != nil {
		b.Meta().SetMAC(mac)
	}

	if topic := component.StateTopic(); topic != "" {
		b.MQTT().Subscribe(mqtt.NewSubscriber(topic, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
			return component.SetState(message)
		}))
	}

	return err
}
