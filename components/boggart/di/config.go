package di

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/config"
)

type ConfigContainerSupport interface {
	SetConfig(*ConfigContainer)
	Config() *ConfigContainer
}

func ConfigContainerBind(bind boggart.Bind) (*ConfigContainer, bool) {
	if support, ok := bind.(ConfigContainerSupport); ok {
		return support.Config(), true
	}

	return nil, false
}

type ConfigBind struct {
	mutex     sync.RWMutex
	container *ConfigContainer
}

func (b *ConfigBind) SetConfig(container *ConfigContainer) {
	b.mutex.Lock()
	b.container = container
	b.mutex.Unlock()
}

func (b *ConfigBind) Config() *ConfigContainer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.container
}

type ConfigContainer struct {
	bind      boggart.BindItem
	configApp config.Component
}

func NewConfigContainer(bind boggart.BindItem, configApp config.Component) *ConfigContainer {
	return &ConfigContainer{
		bind:      bind,
		configApp: configApp,
	}
}

func (b *ConfigContainer) Bind() interface{} {
	return b.bind.Config()
}

func (b *ConfigContainer) App() config.Component {
	return b.configApp
}
