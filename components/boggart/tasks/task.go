package tasks

import (
	"context"
)

/*
Что должны знать о задаче:
  - Сколько времени даем на выполнение / HandlerWithTimeout
  - Сколько раз повторяем / ScheduleWithAttemptsLimit, ScheduleWithSuccessLimit, ScheduleWithFailsLimit
  - Какой интервал повторения (а лучше вычисляемые даты, чтобы крон покрыть), интервал может быть динамическим (например, после ошибки увеличиваем интервал или
    после изменения внутренних значений сразу готовы к выполнению)
  - Иметь возможность завершить даже в процессе выполнения (например чтобы два раза не повторять в один момент времени) / manager.Cancel
  - Иметь возможность запустить задачу внепланово (из интерфейса) / manager.Handle
  - Знать сколько раз задача выполнялась (с разрезами успешно/не успешно) / meta.Attempts, meta.Success, meta.Fails
  - Уметь именовать человекопонятно / task.Name
  - Посмотреть в каком статусе задача (выполняется, ждет выполнения, когда ближайшее выполнение) / meta.* и schedule
*/
type Task interface {
	Name() string
	Handler() Handler
	Schedule() Schedule
}

type TaskBase struct {
	name     string
	handler  Handler
	schedule Schedule
}

func NewTask() *TaskBase {
	return &TaskBase{}
}

func (t *TaskBase) WithName(name string) *TaskBase {
	t.name = name
	return t
}

func (t *TaskBase) WithHandler(handler Handler) *TaskBase {
	t.handler = handler
	return t
}

func (t *TaskBase) WithHandlerFunc(handler func(context.Context) error) *TaskBase {
	return t.WithHandler(HandlerFuncFromShortToLong(handler))
}

func (t *TaskBase) WithHandlerFuncFull(handler HandlerFunc) *TaskBase {
	t.handler = handler
	return t
}

func (t *TaskBase) WithSchedule(schedule Schedule) *TaskBase {
	t.schedule = schedule
	return t
}

func (t *TaskBase) WithScheduleFunc(schedule ScheduleFunc) *TaskBase {
	t.schedule = schedule
	return t
}

func (t *TaskBase) Name() string {
	return t.name
}

func (t *TaskBase) Handler() Handler {
	return t.handler
}

func (t *TaskBase) Schedule() Schedule {
	return t.schedule
}
