package di

import (
	"context"
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	w "github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/shadow/components/workers"
)

type bindTask interface {
	w.Task
	SetName(string)
}

type WorkersHasTasks interface {
	Tasks() []w.Task
}

type WorkersContainerSupport interface {
	SetWorkers(*WorkersContainer)
	Workers() *WorkersContainer
}

func WorkersContainerBind(bind boggart.Bind) (*WorkersContainer, bool) {
	if support, ok := bind.(WorkersContainerSupport); ok {
		container := support.Workers()
		return container, container != nil
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

	client workers.Component

	mutex sync.RWMutex
	tasks []w.Task
}

func NewWorkersContainer(bind boggart.BindItem, client workers.Component) *WorkersContainer {
	return &WorkersContainer{
		bind:   bind,
		client: client,
	}
}

func (c *WorkersContainer) HookRegister() {
	// инициализируем таски из привязки
	if has, ok := c.bind.Bind().(WorkersHasTasks); ok {
		tasks := has.Tasks()
		c.tasks = make([]w.Task, 0, len(tasks)+2) // 2 на пробы запас

		for _, tsk := range has.Tasks() {
			c.RegisterTask(tsk)
		}
	} else {
		c.tasks = make([]w.Task, 0, 2)
	}
}

func (c *WorkersContainer) HookUnregister() {
	for _, tsk := range c.Tasks() {
		c.UnregisterTask(tsk)
	}
}

func (c *WorkersContainer) createTask(tsk w.Task) w.Task {
	if tsk, ok := tsk.(bindTask); ok {
		tsk.SetName("bind-" + c.bind.ID() + "-" + c.bind.Type() + "-" + tsk.Name())
	}

	if logger, ok := LoggerContainerBind(c.bind.Bind()); ok {
		// не логировать дважды ошибку с таски о пробе
		if probes, ok := ProbesContainerBind(c.bind.Bind()); ok {
			if tsk == probes.Liveness() || tsk == probes.Readiness() {
				logger = nil
			}
		}

		tsk = newWorkersWrapTask(tsk, logger)
	}

	return tsk
}

func (c *WorkersContainer) RegisterTask(tsk w.Task) {
	tsk = c.createTask(tsk)

	c.mutex.Lock()
	c.tasks = append(c.tasks, tsk)
	c.mutex.Unlock()

	c.client.AddTask(tsk)
}

func (c *WorkersContainer) UnregisterTask(tsk w.Task) {
	c.client.RemoveTask(tsk)

	c.mutex.Lock()

	for i := len(c.tasks) - 1; i >= 0; i-- {
		if wrap, ok := c.tasks[i].(*workersWrapTask); ok && (wrap.original == tsk || c.tasks[i] == tsk) {
			c.tasks = append(c.tasks[:i], c.tasks[i+1:]...)
		}
	}

	c.mutex.Unlock()
}

func (c *WorkersContainer) Tasks() []w.Task {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return append([]w.Task(nil), c.tasks...)
}

func (c *WorkersContainer) WrapTaskIsOnline(fn func(context.Context) error) *task.FunctionTask {
	return task.NewFunctionTask(func(ctx context.Context) (interface{}, error) {
		if c.bind.Status() != boggart.BindStatusOnline {
			return nil, nil
		}

		return nil, fn(ctx)
	})
}

func (c *WorkersContainer) WrapTaskOnceSuccess(fn func(context.Context) error) (tsk *task.FunctionTask) {
	tsk = task.NewFunctionTask(func(ctx context.Context) (interface{}, error) {
		err := fn(ctx)

		if err == nil {
			c.UnregisterTask(tsk)
		}

		return nil, err
	})

	return tsk
}

func (c *WorkersContainer) WrapTaskIsOnlineOnceSuccess(fn func(context.Context) error) (tsk *task.FunctionTask) {
	tsk = task.NewFunctionTask(func(ctx context.Context) (interface{}, error) {
		if c.bind.Status() != boggart.BindStatusOnline {
			return nil, nil
		}

		err := fn(ctx)

		if err == nil {
			c.UnregisterTask(tsk)

			if attempt, ok := w.AttemptFromContext(ctx); ok {
				tsk.SetRepeats(attempt)
			}
		}

		return nil, err
	})

	return tsk
}

type workersWrapTask struct {
	task.BaseTask

	original w.Task
	logger   *LoggerContainer
}

func newWorkersWrapTask(tsk w.Task, logger *LoggerContainer) *workersWrapTask {
	t := &workersWrapTask{
		original: tsk,
		logger:   logger,
	}
	t.Init()
	t.sync()

	return t
}

// обертки воркеров кривоватые, поэтому хачим синхронизацию состояния
// nolint:golint
func (t *workersWrapTask) Id() string {
	return t.original.Id()
}

func (t *workersWrapTask) sync() {
	t.SetName(t.original.Name())
	t.SetPriority(t.original.Priority())
	t.SetRepeats(t.original.Repeats())
	t.SetRepeatInterval(t.original.RepeatInterval())
	t.SetTimeout(t.original.Timeout())

	if st := t.original.StartedAt(); st != nil {
		t.SetStartedAt(*st)
	}
}

func (t *workersWrapTask) Run(ctx context.Context) (result interface{}, err error) {
	result, err = t.original.Run(ctx)
	if err != nil && t.logger != nil {
		t.logger.Error("Task ended with an error",
			"error", err.Error(),
			"task", t.Name(),
		)
	}

	t.sync()

	return result, err
}
