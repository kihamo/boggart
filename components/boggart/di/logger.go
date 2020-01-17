package di

import (
	"sync"

	"github.com/kihamo/shadow/components/logging"
)

type LoggerContainerSupport interface {
	SetLogger(*LoggerContainer)
	Logger() *LoggerContainer
}

type LoggerBind struct {
	mutex     sync.RWMutex
	container *LoggerContainer
}

func (b *LoggerBind) SetLogger(container *LoggerContainer) {
	b.mutex.Lock()
	b.container = container
	b.mutex.Unlock()
}

func (b *LoggerBind) Logger() *LoggerContainer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.container
}

type LoggerContainer struct {
	logging.Logger
}

func NewLoggerContainer(logger logging.Logger) *LoggerContainer {
	return &LoggerContainer{
		Logger: logger,
	}
}
