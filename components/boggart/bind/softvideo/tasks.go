package softvideo

import (
	"context"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.updaterInterval)
	taskUpdater.SetName("bind-softvideo-updater-" + b.provider.AccountID())

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	value, err := b.provider.Balance(ctx)
	if err != nil {
		b.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.DeviceStatusOnline)

	current := int64(value * 100)
	prev := atomic.LoadInt64(&b.lastValue)

	if current != prev {
		atomic.StoreInt64(&b.lastValue, current)

		b.MQTTPublishAsync(ctx, MQTTTopicBalance.Format(b.SerialNumber()), 0, true, value)
	}

	return nil, nil
}
