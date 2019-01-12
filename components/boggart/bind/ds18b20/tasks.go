package ds18b20

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/yryz/ds18b20"
)

func (b *Bind) Tasks() []workers.Task {
	sn := b.SerialNumber()

	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Minute)
	taskLiveness.SetName("bind-ds18b20-liveness-" + sn)

	taskStateUpdater := task.NewFunctionTask(b.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(time.Minute)
	taskStateUpdater.SetName("bind-ds18b20-state-updater-" + sn)

	return []workers.Task{
		taskLiveness,
		taskStateUpdater,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	devices, err := ds18b20.Sensors()
	if err != nil {
		b.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	sn := b.SerialNumber()

	for _, device := range devices {
		if device == sn {
			b.UpdateStatus(boggart.DeviceStatusOnline)
			return nil, nil
		}
	}

	b.UpdateStatus(boggart.DeviceStatusOffline)
	return nil, nil
}

func (b *Bind) taskStateUpdater(ctx context.Context) (interface{}, error) {
	if b.Status() != boggart.DeviceStatusOnline {
		return nil, nil
	}

	sn := b.SerialNumber()

	value, err := ds18b20.Temperature(sn)
	if err != nil {
		return nil, err
	}

	prev := atomic.LoadInt64(&b.lastValue)
	current := int64(value * 1000)

	if prev != current {
		atomic.StoreInt64(&b.lastValue, current)

		b.MQTTPublishAsync(ctx, MQTTTopicValue.Format(sn), 0, true, value)
	}

	return nil, nil
}
