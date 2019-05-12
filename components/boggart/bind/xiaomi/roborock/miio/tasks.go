package miio

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("bind-xiaomi-roborock-liveness-" + b.config.Host)

	taskState := task.NewFunctionTask(b.taskUpdater)
	taskState.SetTimeout(b.config.UpdaterTimeout)
	taskState.SetRepeats(-1)
	taskState.SetRepeatInterval(b.config.UpdaterInterval)
	taskState.SetName("bind-hikvision-updater-" + b.config.Host)

	return []workers.Task{
		taskLiveness,
		taskState,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	sn, err := b.device.SerialNumber(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, nil
	}

	if b.Status() == boggart.BindStatusOnline {
		return nil, nil
	}

	if b.SerialNumber() == "" {
		b.SetSerialNumber(sn)
	}
	b.UpdateStatus(boggart.BindStatusOnline)

	return nil, nil
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	if b.Status() != boggart.BindStatusOnline {
		return nil, nil
	}

	sn := b.SerialNumber()
	if sn == "" {
		return nil, nil
	}

	snMQTT := mqtt.NameReplace(sn)

	status, err := b.device.Status(ctx)
	if err == nil {
		if ok := b.battery.Set(status.Battery); ok {
			metricBattery.With("serial_number", sn).Set(float64(status.Battery))

			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicBattery.Format(snMQTT), status.Battery); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if ok := b.cleanArea.Set(status.CleanArea); ok {
			metricCleanArea.With("serial_number", sn).Set(float64(status.CleanArea))

			if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicCleanArea.Format(snMQTT), status.CleanArea); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if status.CleanTime > 0 {
			if ok := b.cleanTime.Set(uint32(status.CleanTime)); ok {
				metricCleanTime.With("serial_number", sn).Set(float64(status.CleanTime))

				if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicCleanTime.Format(snMQTT), status.CleanTime); e != nil {
					err = multierr.Append(err, e)
				}
			}
		}
	}

	return nil, err
}
