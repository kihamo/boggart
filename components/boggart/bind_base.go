package boggart

import (
	"context"
	"sync"

	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/shadow/components/logging"
)

type BindBase struct {
	statusManager BindStatusManager
	mutex         sync.RWMutex
	logger        logging.Logger
	id            string
	serialNumber  string
}

func (b *BindBase) Run() error {
	return nil
}

func (b *BindBase) SetID(id string) {
	b.mutex.Lock()
	b.id = id
	b.mutex.Unlock()
}

func (b *BindBase) SetStatusManager(manager BindStatusManager) {
	b.mutex.Lock()
	b.statusManager = manager
	b.mutex.Unlock()
}

func (b *BindBase) Logger() logging.Logger {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.logger == nil {
		b.logger = logging.DefaultLogger()
	}

	return b.logger
}

func (b *BindBase) SetLogger(logger logging.Logger) {
	b.mutex.Lock()
	b.logger = logger
	b.mutex.Unlock()
}

func (b *BindBase) ID() string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.id
}

func (b *BindBase) Status() BindStatus {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.statusManager != nil {
		return b.statusManager()
	}

	return BindStatusUnknown
}

func (b *BindBase) IsStatus(status BindStatus) bool {
	return b.Status() == status
}

func (b *BindBase) IsStatusUnknown() bool {
	return b.IsStatus(BindStatusUnknown)
}

func (b *BindBase) IsStatusUninitialized() bool {
	return b.IsStatus(BindStatusUninitialized)
}

func (b *BindBase) IsStatusInitializing() bool {
	return b.IsStatus(BindStatusInitializing)
}

func (b *BindBase) IsStatusOnline() bool {
	return b.IsStatus(BindStatusOnline)
}

func (b *BindBase) IsStatusOffline() bool {
	return b.IsStatus(BindStatusOffline)
}

func (b *BindBase) IsStatusRemoving() bool {
	return b.IsStatus(BindStatusRemoving)
}

func (b *BindBase) IsStatusRemoved() bool {
	return b.IsStatus(BindStatusRemoved)
}

func (b *BindBase) WrapTaskIsOnline(fn func(context.Context) error) *task.FunctionTask {
	return task.NewFunctionTask(func(ctx context.Context) (interface{}, error) {
		if !b.IsStatusOnline() {
			return nil, nil
		}

		return nil, fn(ctx)
	})
}

func (b *BindBase) SerialNumber() string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.serialNumber
}

func (b *BindBase) SetSerialNumber(serialNumber string) {
	b.mutex.Lock()
	b.serialNumber = serialNumber
	b.mutex.Unlock()
}
