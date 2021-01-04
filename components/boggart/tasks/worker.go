package tasks

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/kihamo/boggart/atomic"
)

// внутренняя сущность, которая используется в менеджере для работы с задачей
type worker struct {
	io.Closer

	readyFn  func(string)
	removeFn func(string)

	task Task
	meta *Meta

	cancelMutex sync.Mutex
	cancelFn    context.CancelFunc

	schedulerOnce   atomic.Once
	schedulerTicker *time.Ticker
	schedulerChange chan time.Duration
	schedulerDone   chan struct{}
}

func newWorker(ready, remove func(string), task Task) *worker {
	return &worker{
		readyFn:         ready,
		removeFn:        remove,
		task:            task,
		meta:            newMeta(),
		schedulerChange: make(chan time.Duration, 1),
		schedulerDone:   make(chan struct{}),
	}
}

func (w *worker) RunScheduler() {
	now := time.Now()

	next := w.task.Schedule().Next(*w.meta)
	if next.IsZero() {
		// fmt.Println("Remove IsZero", w.meta.id, w.task.Name())

		// FIXME: надо различать случай, когда время не возможно определить из-за не сработавшего фактора
		// например, объект в данный момент не активен и поэтому задача для него не выполняется, но как только
		// он станет активен, задачу можно выполнить
		w.removeFn(w.meta.id)
		return
	}

	if now.After(next) {
		// fmt.Println("Remove After", w.meta.id, w.task.Name())

		w.removeFn(w.meta.id)
		return
	}

	w.meta.nextRunAt.Set(next)

	duration := next.Sub(now)

	// FIXME: dirty hack, но пока так, точность запуска - 1 секунда
	// между time.Now() и schedule.Next() который так же возвращает time.Now() проходит время в несколько наносекунд
	// и этот оверхед дает дельту, хотя фактически имелось ввиду одно и тоже время, поэтому принудительно сфильтровываем
	// все что меньше секунды
	if duration < time.Second {
		go w.readyFn(w.meta.id)
		return
	}

	needChangeDuration := true

	// тикер запускаем только один раз
	w.schedulerOnce.Do(func() {
		needChangeDuration = false
		w.schedulerTicker = time.NewTicker(duration)

		go func() {
			defer w.schedulerTicker.Stop()

			for {
				select {
				case <-w.schedulerTicker.C:
					if s := w.meta.Status(); s.IsAllowExecute() {
						w.readyFn(w.meta.id)
					}

				case d := <-w.schedulerChange:
					w.schedulerTicker.Reset(d)

				case <-w.schedulerDone:
					return
				}
			}
		}()
	})

	if needChangeDuration {
		w.schedulerChange <- duration
	}
}

func (w *worker) Handle(ctx context.Context) error {
	if s := w.meta.Status(); !s.IsAllowExecute() {
		return ErrAlreadyRunning
	}

	now := time.Now()
	if w.meta.firstRunAt.IsNil() {
		w.meta.firstRunAt.Set(now)
	}
	w.meta.lastRunAt.Set(now)

	w.meta.attempts.Inc()

	w.meta.status.Set(StatusRunning.Uint32())
	defer w.meta.status.Set(StatusWaiting.Uint32())
	defer w.RunScheduler()

	// оставляем себе возможность принудительно прервать выполнение задачи
	w.cancelMutex.Lock()
	ctx, w.cancelFn = context.WithCancel(ctx)
	w.cancelMutex.Unlock()

	metricHandleStatus.With("status", "started").Inc()

	err := w.task.Handler().Handle(ctx, *w.meta, w.task)

	metricHandleStatus.With("status", "finished").Inc()
	if err == nil {
		metricHandleStatus.With("status", "success").Inc()
	} else {
		metricHandleStatus.With("status", "fail").Inc()
	}

	w.cancelMutex.Lock()
	w.cancelFn = nil
	w.cancelMutex.Unlock()

	if err != nil {
		w.meta.fails.Inc()
	}

	return err
}

func (w *worker) Cancel() {
	w.cancelMutex.Lock()
	cancel := w.cancelFn
	w.cancelMutex.Unlock()

	if cancel != nil {
		cancel()
	}
}

func (w *worker) Close() error {
	w.Cancel()
	w.meta.status.Set(StatusClosed.Uint32())

	close(w.schedulerDone)

	return nil
}
