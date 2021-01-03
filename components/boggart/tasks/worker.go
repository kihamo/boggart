package tasks

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/pborman/uuid"
)

// внутренняя сущность, которая используется в менеджере для работы с задачей
type worker struct {
	io.Closer

	readyFn  func(string)
	removeFn func(string)

	task Task
	meta Meta

	cancelMutex sync.Mutex
	cancelFn    context.CancelFunc

	schedulerOnce   atomic.Once
	schedulerTicker *time.Ticker
	schedulerChange chan time.Duration
	schedulerDone   chan struct{}
}

func newWorker(ready, remove func(string), task Task) *worker {
	w := &worker{
		readyFn:         ready,
		removeFn:        remove,
		task:            task,
		schedulerChange: make(chan time.Duration, 1),
		schedulerDone:   make(chan struct{}),
	}
	w.meta.id = uuid.New()

	return w
}

func (w *worker) RunScheduler() {
	now := time.Now()

	next := w.task.Schedule().Next(w.meta)
	if next.IsZero() {
		// FIXME: надо различать случай, когда время не возможно определить из-за не сработавшего фактора
		// например, объект в данный момент не активен и поэтому задача для него не выполняется, но как только
		// он станет активен, задачу можно выполнить
		w.removeFn(w.meta.id)
		return
	}

	if now.After(next) {
		w.removeFn(w.meta.id)
		return
	}

	duration := next.Sub(now)

	// тикер запускаем только один раз
	w.schedulerOnce.Do(func() {
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

	w.schedulerChange <- duration
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

	err := w.task.Handler().Handle(ctx, w.meta, w.task)

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
