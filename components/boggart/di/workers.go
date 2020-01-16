package di

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
)

type WorkersHasTasks interface {
	Tasks() []workers.Task
}

type WorkersHasListeners interface {
	Listeners() []workers.ListenerWithEvents
}

type WorkersContainerSupport interface {
	SetWorkersContainer(*WorkersContainer)
	WorkersContainer() *WorkersContainer
}

type WorkersBind struct {
	mutex     sync.RWMutex
	container *WorkersContainer
}

func (b *WorkersBind) SetWorkersContainer(container *WorkersContainer) {
	b.mutex.Lock()
	b.container = container
	b.mutex.Unlock()
}

func (b *WorkersBind) WorkersContainer() *WorkersContainer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.container
}

type WorkersContainer struct {
	bind boggart.Bind

	cacheMutex     sync.Mutex
	cacheTasks     []workers.Task
	cacheListeners []workers.ListenerWithEvents
}

func NewWorkersContainer(bind boggart.Bind) *WorkersContainer {
	return &WorkersContainer{
		bind: bind,
	}
}

func (c *WorkersContainer) Listeners() []workers.ListenerWithEvents {
	has, ok := c.bind.(WorkersHasListeners)
	if !ok {
		return nil
	}

	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	if c.cacheListeners == nil {
		c.cacheListeners = has.Listeners()
	}

	return c.cacheListeners
}

func (c *WorkersContainer) Tasks() []workers.Task {
	has, ok := c.bind.(WorkersHasTasks)
	if !ok {
		return nil
	}

	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	if c.cacheTasks == nil {
		c.cacheTasks = has.Tasks()
	}

	return c.cacheTasks
}
