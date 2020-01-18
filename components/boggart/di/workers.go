package di

import (
	"context"
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/shadow/components/logging"
)

type bindTask interface {
	workers.Task
	SetName(string)
}

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

func WorkersContainerBind(bind boggart.Bind) (*WorkersContainer, bool) {
	if support, ok := bind.(WorkersContainerSupport); ok {
		return support.Workers(), true
	}

	return nil, false
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
		tasks := has.Tasks()

		if bindSupport, ok := c.bind.Bind().(LoggerContainerSupport); ok {
			logger := bindSupport.Logger()

			for i, tsk := range tasks {
				if tsk, ok := tsk.(bindTask); ok {
					tsk.SetName("bind-" + c.bind.ID() + "-" + c.bind.Type() + "-" + tsk.Name())
				}

				tasks[i] = c.WrapTaskLogError(tsk, logger)
			}
		}

		c.cacheTasks = tasks
	}

	return c.cacheTasks
}

func (c *WorkersContainer) WrapTaskLogError(t workers.Task, logger logging.Logger) workers.Task {
	n := task.NewFunctionTask(func(ctx context.Context) (result interface{}, err error) {
		result, err = t.Run(ctx)
		if err != nil {
			logger.Error("Task ended with an error",
				"error", err.Error(),
				"task", t.Name(),
			)
		}

		return result, err
	})

	n.SetName(t.Name())
	n.SetPriority(t.Priority())
	n.SetRepeats(t.Repeats())
	n.SetRepeatInterval(t.RepeatInterval())
	n.SetTimeout(t.Timeout())

	if st := t.StartedAt(); st != nil {
		n.SetStartedAt(*st)
	}

	return n
}

func (c *WorkersContainer) WrapTaskIsOnline(fn func(context.Context) error) *task.FunctionTask {
	return task.NewFunctionTask(func(ctx context.Context) (interface{}, error) {
		if c.bind.Status() != boggart.BindStatusOnline {
			return nil, nil
		}

		return nil, fn(ctx)
	})
}
