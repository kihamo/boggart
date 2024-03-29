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
		container := support.Config()
		return container, container != nil
	}

	return nil, false
}

func ConfigForBind(bind boggart.Bind) (interface{}, bool) {
	if support, ok := ConfigContainerBind(bind); ok {
		return support.Bind(), true
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
	configBind interface{}
	configApp  config.Component
}

func NewConfigContainer(configBind interface{}, configApp config.Component) *ConfigContainer {
	return &ConfigContainer{
		configBind: configBind,
		configApp:  configApp,
	}
}

func (b *ConfigContainer) Bind() interface{} {
	return b.configBind
}

func (b *ConfigContainer) App() config.Component {
	return b.configApp
}
