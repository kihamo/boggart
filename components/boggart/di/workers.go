package di

import (
	"context"
	"sync"

	"github.com/kihamo/go-workers/task"

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
	SetWorkers(*WorkersContainer)
	Workers() *WorkersContainer
}

type WorkersBind struct {
	mutex     sync.RWMutex
	container *WorkersContainer
}

func (b *WorkersBind) SetWorkers(container *WorkersContainer) {
	b.mutex.Lock()
	b.container = container
	b.mutex.Unlock()
}

func (b *WorkersBind) Workers() *WorkersContainer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.container
}

type WorkersContainer struct {
	bind boggart.BindItem

	cacheMutex     sync.Mutex
	cacheTasks     []workers.Task
	cacheListeners []workers.ListenerWithEvents
}

func NewWorkersContainer(bind boggart.BindItem) *WorkersContainer {
	return &WorkersContainer{
		bind: bind,
	}
}

func (c *WorkersContainer) Listeners() []workers.ListenerWithEvents {
	has, ok := c.bind.Bind().(WorkersHasListeners)
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
	has, ok := c.bind.Bind().(WorkersHasTasks)
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

func (c *WorkersContainer) WrapTaskIsOnline(fn func(context.Context) error) *task.FunctionTask {
	return task.NewFunctionTask(func(ctx context.Context) (interface{}, error) {
		if c.bind.Status() != boggart.BindStatusOnline {
			return nil, nil
		}

		return nil, fn(ctx)
	})
}
