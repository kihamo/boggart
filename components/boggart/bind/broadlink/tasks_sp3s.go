package broadlink

import (
	"context"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	er "github.com/pkg/errors"
	"go.uber.org/multierr"
)

func (b *BindSP3S) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.updaterInterval)
	taskUpdater.SetName("bind-broadlink:sp3s-updater-" + b.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (b *BindSP3S) taskUpdater(ctx context.Context) (interface{}, error) {
	state, err := b.State()
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	serialNumber := b.SerialNumber()
	serialNumberMQTT := mqtt.NameReplace(serialNumber)
	var result error

	prevState := atomic.LoadInt64(&b.state)
	if prevState == 0 || (prevState == 1) != state {
		if state {
			atomic.StoreInt64(&b.state, 1)
		} else {
			atomic.StoreInt64(&b.state, -1)
		}

		if err := b.MQTTPublishAsync(ctx, SP3SMQTTPublishTopicState.Format(serialNumberMQTT), 0, true, state); err != nil {
			result = multierr.Append(result, err)
		}
	}

	value, err := b.Power()
	if err == nil {
		currentPower := int64(value * 100)
		prevPower := atomic.LoadInt64(&b.power)

		if currentPower != prevPower {
			atomic.StoreInt64(&b.power, currentPower)
			metricSP3SPower.With("serial_number", serialNumber).Set(value)

			if err := b.MQTTPublishAsync(ctx, SP3SMQTTPublishTopicPower.Format(serialNumberMQTT), 0, true, value); err != nil {
				result = multierr.Append(result, err)
			}
		}
	}

	if result != nil {
		result = er.Wrap(result, "Failed send to MQTT")
	}

	return nil, result
}