package sp3s

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(b.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.updaterInterval)
	taskUpdater.SetName("bind-broadlink:sp3s-updater-" + b.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	state, err := b.State()
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	serialNumber := b.SerialNumber()
	serialNumberMQTT := mqtt.NameReplace(serialNumber)

	if ok := b.state.Set(state); ok {
		if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicState.Format(serialNumberMQTT), state); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if power, e := b.Power(); e == nil {
		if ok := b.power.Set(float32(power)); ok {
			metricPower.With("serial_number", serialNumber).Set(power)

			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicPower.Format(serialNumberMQTT), power); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	return nil, err
}
