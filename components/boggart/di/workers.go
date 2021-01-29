package di

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/tasks"
)

type BindHasTasks interface {
	Tasks() []tasks.Task
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

	manager *tasks.Manager

	tasksIDMutex sync.RWMutex
	tasksID      [][2]string
}

func NewWorkersContainer(bind boggart.BindItem, manager *tasks.Manager) *WorkersContainer {
	return &WorkersContainer{
		bind:    bind,
		manager: manager,
		tasksID: make([][2]string, 0, 2),
	}
}

func (c *WorkersContainer) HookRegister() (err error) {
	has, ok := c.bind.Bind().(BindHasTasks)
	if !ok {
		return nil
	}

	list := has.Tasks()
	tmpList := make([][2]string, 0, len(list)+2) // 2 на пробы запас

	var id string

	for _, tsk := range list {
		id, err = c.RegisterTask(tsk)
		if err != nil {
			break
		}

		tmpList = append(tmpList, [2]string{id, tsk.Name()})
	}

	c.tasksIDMutex.Lock()
	c.tasksID = tmpList
	c.tasksIDMutex.Unlock()

	return err
}

func (c *WorkersContainer) HookUnregister() {
	c.tasksIDMutex.Lock()
	defer c.tasksIDMutex.Unlock()

	for _, tsk := range c.tasksID {
		c.manager.Unregister(tsk[0])
	}

	c.tasksID = c.tasksID[:0]
}

func (c *WorkersContainer) TaskShortName(name string) string {
	return strings.TrimPrefix(name, c.prefixTaskName())
}

func (c *WorkersContainer) prefixTaskName() string {
	return "bind/" + c.bind.Type() + "/" + c.bind.ID() + "/"
}

func (c *WorkersContainer) RegisterTask(tsk tasks.Task) (id string, err error) {
	if t, ok := tsk.(*tasks.TaskBase); ok {
		if logger, ok := LoggerContainerBind(c.bind.Bind()); ok {
			// не логировать дважды ошибку с таски о пробе
			if _, ok := ProbesContainerBind(c.bind.Bind()); ok {
				shortName := c.TaskShortName(t.Name())

				if shortName == ProbesConfigReadinessDefaultName || shortName == ProbesConfigLivenessDefaultName {
					logger = nil
				}
			}

			if logger != nil {
				t.WithHandler(tasks.HandlerWithLogger(t.Handler(), logger))
			}
		}

		tsk = t.WithName(c.prefixTaskName() + t.Name())
	}

	id, err = c.manager.Register(tsk)
	if err == nil {
		c.tasksIDMutex.Lock()
		c.tasksID = append(c.tasksID, [2]string{id, tsk.Name()})
		c.tasksIDMutex.Unlock()
	}

	return id, err
}

func (c *WorkersContainer) TasksID() [][2]string {
	c.tasksIDMutex.RLock()
	defer c.tasksIDMutex.RUnlock()

	return append([][2]string(nil), c.tasksID...)
}

func (c *WorkersContainer) TaskInfoByID(id string) (tasks.Task, *tasks.Meta, error) {
	task, err := c.manager.Task(id)
	if err != nil {
		return nil, nil, err
	}

	meta, err := c.manager.Meta(id)
	if err != nil {
		return nil, nil, err
	}

	return task, &meta, err
}

func (c *WorkersContainer) TaskRunByID(ctx context.Context, id string) error {
	return c.manager.Handle(ctx, id)
}

func (c *WorkersContainer) TaskRunByName(ctx context.Context, name string) error {
	name = c.prefixTaskName() + name

	for _, item := range c.TasksID() {
		if item[1] == name {
			if err := c.TaskRunByID(ctx, item[0]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *WorkersContainer) ScheduleRecalculateByID(id string) error {
	return c.manager.Recalculate(id)
}

func (c *WorkersContainer) ScheduleRecalculateByName(name string) error {
	name = c.prefixTaskName() + name

	for _, item := range c.TasksID() {
		if item[1] == name {
			if err := c.ScheduleRecalculateByID(item[0]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *WorkersContainer) WrapTaskHandlerIsOnline(parent tasks.Handler) tasks.Handler {
	return tasks.HandlerFunc(func(ctx context.Context, meta tasks.Meta, task tasks.Task) error {
		if parent == nil {
			return tasks.ErrParentHandlerIsNil
		}

		if !c.bind.Status().IsStatusOnline() {
			return errors.New("bind isn't online")
		}

		return parent.Handle(ctx, meta, task)
	})
}
