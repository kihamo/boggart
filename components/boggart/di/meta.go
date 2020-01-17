package di

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
)

type MetaContainerSupport interface {
	SetMeta(*MetaContainer)
	Meta() *MetaContainer
}

func MetaContainerBind(bind boggart.Bind) (*MetaContainer, bool) {
	if support, ok := bind.(MetaContainerSupport); ok {
		return support.Meta(), true
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
	bind boggart.BindItem

	mutex        sync.RWMutex
	serialNumber string
}

func NewMetaContainer(bind boggart.BindItem) *MetaContainer {
	return &MetaContainer{
		bind: bind,
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

func (b *MetaContainer) IsStatus(status boggart.BindStatus) bool {
	return b.Status() == status
}

func (b *MetaContainer) IsStatusUnknown() bool {
	return b.IsStatus(boggart.BindStatusUnknown)
}

func (b *MetaContainer) IsStatusUninitialized() bool {
	return b.IsStatus(boggart.BindStatusUninitialized)
}

func (b *MetaContainer) IsStatusInitializing() bool {
	return b.IsStatus(boggart.BindStatusInitializing)
}

func (b *MetaContainer) IsStatusOnline() bool {
	return b.IsStatus(boggart.BindStatusOnline)
}

func (b *MetaContainer) IsStatusOffline() bool {
	return b.IsStatus(boggart.BindStatusOffline)
}

func (b *MetaContainer) IsStatusRemoving() bool {
	return b.IsStatus(boggart.BindStatusRemoving)
}

func (b *MetaContainer) IsStatusRemoved() bool {
	return b.IsStatus(boggart.BindStatusRemoved)
}

func (b *MetaContainer) SerialNumber() string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.serialNumber
}

func (b *MetaContainer) SetSerialNumber(serialNumber string) {
	b.mutex.Lock()
	b.serialNumber = serialNumber
	b.mutex.Unlock()
}
