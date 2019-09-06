package ds18b20

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/yryz/ds18b20"
)

func (b *Bind) Tasks() []workers.Task {
	sn := b.SerialNumber()

	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.livenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.livenessInterval)
	taskLiveness.SetName("liveness-" + sn)

	taskStateUpdater := task.NewFunctionTask(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.updaterInterval)
	taskStateUpdater.SetName("updater-" + sn)

	return []workers.Task{
		taskLiveness,
		taskStateUpdater,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	devices, err := ds18b20.Sensors()
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	sn := b.SerialNumber()

	for _, device := range devices {
		if device == sn {
			b.UpdateStatus(boggart.BindStatusOnline)
			return nil, nil
		}
	}

	b.UpdateStatus(boggart.BindStatusOffline)
	return nil, nil
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	if b.Status() != boggart.BindStatusOnline {
		return nil, nil
	}

	value, err := b.Temperature()
	if err != nil {
		return nil, err
	}

	sn := b.SerialNumber()

	metricValue.With("serial_number", sn).Set(value)

	if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicValue.Format(mqtt.NameReplace(sn)), value); err != nil {
		return nil, err
	}

	return nil, nil
}
