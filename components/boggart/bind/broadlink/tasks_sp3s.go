package broadlink

import (
	"context"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *BindSP3S) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskStateUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.updaterInterval)
	taskUpdater.SetName("bind-broadlink-sp3s-updater-" + b.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (b *BindSP3S) taskStateUpdater(ctx context.Context) (interface{}, error) {
	state, err := b.State()
	if err != nil {
		b.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.DeviceStatusOnline)

	serialNumber := b.SerialNumber()
	serialNumberMQTT := mqtt.NameReplace(serialNumber)

	prevState := atomic.LoadInt64(&b.state)
	if prevState == 0 || (prevState == 1) != state {
		if state {
			atomic.StoreInt64(&b.state, 1)
		} else {
			atomic.StoreInt64(&b.state, -1)
		}

		b.MQTTPublishAsync(ctx, SP3SMQTTTopicState.Format(serialNumberMQTT), 0, true, state)
	}

	value, err := b.Power()
	if err != nil {
		return nil, nil
	}

	metricSP3SPower.With("serial_number", serialNumber).Set(value)

	currentPower := int64(value * 100)
	prevPower := atomic.LoadInt64(&b.power)

	if currentPower != prevPower {
		atomic.StoreInt64(&b.power, currentPower)

		b.MQTTPublishAsync(ctx, SP3SMQTTTopicPower.Format(serialNumberMQTT), 0, true, value)
	}

	return nil, nil
}
