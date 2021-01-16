package di

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/config"
)

type MetaContainerSupport interface {
	SetMeta(*MetaContainer)
	Meta() *MetaContainer
}

func MetaContainerBind(bind boggart.Bind) (*MetaContainer, bool) {
	if support, ok := bind.(MetaContainerSupport); ok {
		container := support.Meta()
		return container, container != nil
	}

	return nil, false
}

type MetaBind struct {
	mutex     sync.RWMutex
	container *MetaContainer
}

func (b *MetaBind) SetMeta(container *MetaContainer) {
	b.mutex.Lock()
	b.container = container
	b.mutex.Unlock()
}

func (b *MetaBind) Meta() *MetaContainer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.container
}

type MetaContainer struct {
	bind   boggart.BindItem
	mqtt   mqtt.Component
	config config.Component

	serialNumber atomic.Value
	mac          atomic.Value
}

func NewMetaContainer(bind boggart.BindItem, mqtt mqtt.Component, config config.Component) *MetaContainer {
	return &MetaContainer{
		bind:   bind,
		mqtt:   mqtt,
		config: config,
	}
}

func (b *MetaContainer) BindType() boggart.BindType {
	return b.bind.BindType()
}

func (b *MetaContainer) ID() string {
	return b.bind.ID()
}

func (b *MetaContainer) Type() string {
	return b.bind.Type()
}

func (b *MetaContainer) Description() string {
	return b.bind.Description()
}

func (b *MetaContainer) Tags() []string {
	return b.bind.Tags()
}

func (b *MetaContainer) Config() interface{} {
	return b.bind.Config()
}

func (b *MetaContainer) Status() boggart.BindStatus {
	return b.bind.Status()
}

func (b *MetaContainer) SerialNumber() string {
	if sn := b.serialNumber.Load(); sn != nil {
		return sn.(string)
	}

	return ""
}

func (b *MetaContainer) SetSerialNumber(serialNumber string) {
	b.serialNumber.Store(serialNumber)

	b.mqtt.PublishAsyncWithCache(context.Background(), b.MQTTTopicSerialNumber(), 1, true, serialNumber)
}

func (b *MetaContainer) MAC() *net.HardwareAddr {
	if mac := b.mac.Load(); mac != nil {
		return mac.(*net.HardwareAddr)
	}

	return nil
}

func (b *MetaContainer) MACAsString() string {
	if mac := b.MAC(); mac != nil {
		return mac.String()
	}

	return ""
}

func (b *MetaContainer) SetMAC(mac net.HardwareAddr) {
	b.mac.Store(&mac)

	b.mqtt.PublishAsyncWithCache(context.Background(), b.MQTTTopicMAC(), 1, true, mac)
}

func (b *MetaContainer) SetMACAsString(mac string) error {
	addr, err := net.ParseMAC(mac)
	if err == nil {
		b.SetMAC(addr)
	}

	return err
}

func (b *MetaContainer) MQTTTopicStatus() mqtt.Topic {
	return mqtt.Topic(b.config.String(boggart.ConfigMQTTTopicBindStatus)).Format(b.ID())
}

func (b *MetaContainer) MQTTTopicSerialNumber() mqtt.Topic {
	return mqtt.Topic(b.config.String(boggart.ConfigMQTTTopicBindSerialNumber)).Format(b.ID())
}

func (b *MetaContainer) MQTTTopicMAC() mqtt.Topic {
	return mqtt.Topic(b.config.String(boggart.ConfigMQTTTopicBindMAC)).Format(b.ID())
}
