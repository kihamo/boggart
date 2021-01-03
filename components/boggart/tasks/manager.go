package tasks

import (
	"context"
	"io"
	"sync"
)

type Manager struct {
	io.Closer

	workers      map[string]*worker
	workersMutex sync.RWMutex

	ready  chan string
	remove chan string
	done   chan struct{}
}

func NewManager() *Manager {
	m := &Manager{
		workers: make(map[string]*worker),
		ready:   make(chan string, 1), // буферизированный чтобы двойного выполнения не было за время пока идет обработка другого задания
		remove:  make(chan string),
		done:    make(chan struct{}),
	}
	go m.execute()

	return m
}

func (m *Manager) execute() {
	ctx := context.Background()

	// TODO: rate limiter

	for {
		select {
		case id := <-m.ready:
			if w := m.worker(id); w != nil {
				go w.Handle(ctx)
			}

		case id := <-m.remove:
			m.Unregister(id)

		case <-m.done:
			return
		}
	}
}

func (m *Manager) taskReady(id string) {
	m.ready <- id
}

func (m *Manager) taskRemove(id string) {
	m.remove <- id
}

func (m *Manager) worker(id string) *worker {
	m.workersMutex.RLock()
	defer m.workersMutex.RUnlock()

	return m.workers[id]
}

func (m *Manager) Register(task Task) (id string, _ error) {
	if task == nil {
		return id, ErrTaskIsEmpty
	}

	if task.Handler() == nil {
		return id, ErrHandlerIsEmpty
	}

	if task.Schedule() == nil {
		return id, ErrScheduleIsEmpty
	}

	worker := newWorker(m.taskReady, m.taskRemove, task)

	m.workersMutex.Lock()
	m.workers[worker.meta.id] = worker
	m.workersMutex.Unlock()

	worker.RunScheduler()

	return worker.meta.id, nil
}

func (m *Manager) Unregister(id string) {
	if w := m.worker(id); w != nil {
		w.Close()

		m.workersMutex.Lock()
		delete(m.workers, id)
		m.workersMutex.Unlock()
	}
}

func (m *Manager) Handle(ctx context.Context, id string) error {
	if w := m.worker(id); w != nil {
		return w.Handle(ctx)
	}

	return ErrTaskNotFound
}

func (m *Manager) Schedule(id string) (schedule Schedule, _ error) {
	if w := m.worker(id); w != nil {
		return w.task.Schedule(), nil
	}

	return schedule, ErrTaskNotFound
}

func (m *Manager) Meta(id string) (meta Meta, _ error) {
	if w := m.worker(id); w != nil {
		return *w.meta, nil
	}

	return meta, ErrTaskNotFound
}

func (m *Manager) Task(id string) (task Task, _ error) {
	if w := m.worker(id); w != nil {
		return w.task, nil
	}

	return nil, ErrTaskNotFound
}

func (m *Manager) Cancel(id string) {
	if w := m.worker(id); w != nil {
		w.Cancel()
	}
}

// Принудительно обновляет внутренний шедулер в воркере, на случай если произошли изменений
func (m *Manager) Recalculate(id string) error {
	if w := m.worker(id); w != nil {
		w.RunScheduler()
		return nil
	}

	return ErrTaskNotFound
}

func (m *Manager) Close() (err error) {
	m.workersMutex.Lock()
	defer m.workersMutex.Unlock()

	for _, worker := range m.workers {
		if err = worker.Close(); err != nil {
			return err
		}
	}

	close(m.done)
	close(m.ready)
	close(m.remove)

	return err
}
