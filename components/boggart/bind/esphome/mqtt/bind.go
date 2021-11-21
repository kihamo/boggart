package mqtt

import (
	"net"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	components sync.Map
	status     atomic.BoolNull
	ip         *atomic.Value
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.status.Nil()

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

func (b *Bind) delete(id string) {
	if item, ok := b.components.Load(id); ok {
		if cmp, ok := item.(Component); ok {
			b.MQTT().Unsubscribe(cmp.Subscribers()...)
		}
	}

	b.components.Delete(id)
}

func (b *Bind) register(component Component) (err error) {
	// overwrite component
	if _, ok := b.components.Load(component.UniqueID()); ok {
		b.delete(component.UniqueID())
	}

	b.components.Store(component.UniqueID(), component)

	if mac := component.DeviceInfo().MAC(); mac != nil {
		b.Meta().SetMAC(mac)
	}

	return b.MQTT().Subscribe(component.Subscribers()...)
}
