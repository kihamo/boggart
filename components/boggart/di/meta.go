package di

import (
	"context"
	"net"
	"net/url"
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
	bindItem boggart.BindItem
	mqtt     mqtt.Component
	config   config.Component
	diMQTT   *MQTTContainer

	serialNumber atomic.Value
	mac          atomic.Value
	link         atomic.Value
}

func NewMetaContainer(bindItem boggart.BindItem, mqtt mqtt.Component, config config.Component) *MetaContainer {
	ctr := &MetaContainer{
		bindItem: bindItem,
		mqtt:     mqtt,
		config:   config,
	}

	if ctrMQTT, ok := MQTTContainerBind(bindItem.Bind()); ok {
		ctr.diMQTT = ctrMQTT
	}

	return ctr
}

func (b *MetaContainer) BindType() boggart.BindType {
	return b.bindItem.BindType()
}

func (b *MetaContainer) ID() string {
	return b.bindItem.ID()
}

func (b *MetaContainer) Type() string {
	return b.bindItem.Type()
}

func (b *MetaContainer) Description() string {
	return b.bindItem.Description()
}

func (b *MetaContainer) Tags() []string {
	return b.bindItem.Tags()
}

func (b *MetaContainer) Status() boggart.BindStatus {
	return b.bindItem.Status()
}

func (b *MetaContainer) SerialNumber() string {
	if sn := b.serialNumber.Load(); sn != nil {
		return sn.(string)
	}

	return ""
}

func (b *MetaContainer) SetSerialNumber(serialNumber string) {
	b.serialNumber.Store(serialNumber)

	if b.diMQTT != nil {
		b.diMQTT.PublishAsync(context.Background(), b.MQTTTopicSerialNumber(), serialNumber)
		return
	}

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

	if b.diMQTT != nil {
		b.diMQTT.PublishAsync(context.Background(), b.MQTTTopicMAC(), mac)
		return
	}

	b.mqtt.PublishAsyncWithCache(context.Background(), b.MQTTTopicMAC(), 1, true, mac)
}

func (b *MetaContainer) SetMACAsString(mac string) error {
	addr, err := net.ParseMAC(mac)
	if err == nil {
		b.SetMAC(addr)
	}

	return err
}

func (b *MetaContainer) Link() *url.URL {
	if sn := b.link.Load(); sn != nil {
		if link, ok := sn.(*url.URL); ok {
			copy := &url.URL{}
			*copy = *link
			return copy
		}
	}

	return nil
}

func (b *MetaContainer) SetLink(link *url.URL) {
	if link != nil {
		copy := &url.URL{}
		*copy = *link
		link = copy
	}

	b.link.Store(link)

	if b.diMQTT != nil {
		b.diMQTT.PublishAsync(context.Background(), b.MQTTTopicLink(), link)
		return
	}

	b.mqtt.PublishAsyncWithCache(context.Background(), b.MQTTTopicLink(), 1, true, link)
}

func (b *MetaContainer) MQTTTopicStatus() mqtt.Topic {
	return mqtt.Topic(b.config.String(boggart.ConfigMQTTTopicBindStatus)).Format(b.ID())
}

func (b *MetaContainer) MQTTTopicReload() mqtt.Topic {
	return mqtt.Topic(b.config.String(boggart.ConfigMQTTTopicBindReload)).Format(b.ID())
}

func (b *MetaContainer) MQTTTopicSerialNumber() mqtt.Topic {
	return mqtt.Topic(b.config.String(boggart.ConfigMQTTTopicBindSerialNumber)).Format(b.ID())
}

func (b *MetaContainer) MQTTTopicMAC() mqtt.Topic {
	return mqtt.Topic(b.config.String(boggart.ConfigMQTTTopicBindMAC)).Format(b.ID())
}

func (b *MetaContainer) MQTTTopicLink() mqtt.Topic {
	return mqtt.Topic(b.config.String(boggart.ConfigMQTTTopicBindLink)).Format(b.ID())
}
