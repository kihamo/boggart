package softvideo

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
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
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.BindStatusOnline)
	current := float32(value)

	if ok := b.balance.Set(current); ok {
		sn := b.SerialNumber()
		metricBalance.With("account", sn).Set(value)
		if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicBalance.Format(mqtt.NameReplace(sn)), value); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
