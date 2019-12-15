package miio

import (
	"context"
	"errors"
	"time"

	"github.com/kihamo/boggart/providers/xiaomi/miio/devices/vacuum"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillSuccessTask(b.taskSerialNumber)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Second * 30)
	taskSerialNumber.SetName("serial-number")

	taskState := b.WrapTaskIsOnline(b.taskUpdater)
	taskState.SetTimeout(b.config.UpdaterTimeout)
	taskState.SetRepeats(-1)
	taskState.SetRepeatInterval(b.config.UpdaterInterval)
	taskState.SetName("updater")

	return []workers.Task{
		taskSerialNumber,
		taskState,
	}
}

func (b *Bind) taskSerialNumber(ctx context.Context) (interface{}, error) {
	if !b.IsStatusOnline() {
		return nil, errors.New("bind isn't online")
	}

	sn, err := b.device.SerialNumber(ctx)
	if err != nil {
		return nil, err
	}

	b.SetSerialNumber(sn)

	return nil, nil
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	sn := b.SerialNumber()
	if sn == "" {
		return nil
	}

	var err error

	// only statistics
	status, e := b.device.Status(ctx)
	if e == nil {
		metricBattery.With("serial_number", sn).Set(float64(status.Battery))

		if e := b.MQTTPublishAsync(ctx, b.config.TopicBattery.Format(sn), status.Battery); e != nil {
			err = multierr.Append(err, e)
		}

		if e := b.MQTTPublishAsync(ctx, b.config.TopicCleanArea.Format(sn), status.CleanArea); e != nil {
			err = multierr.Append(err, e)
		}

		if status.CleanTime > 0 {
			if e := b.MQTTPublishAsync(ctx, b.config.TopicCleanTime.Format(sn), status.CleanTime); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	consumables, e := b.device.Consumables(ctx)
	if e == nil {
		if consumable, ok := consumables[vacuum.ConsumableFilter]; ok {
			if e := b.MQTTPublishAsync(ctx, b.config.TopicConsumableFilter.Format(sn), consumable); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if consumable, ok := consumables[vacuum.ConsumableBrushMain]; ok {
			if e := b.MQTTPublishAsync(ctx, b.config.TopicConsumableBrushMain.Format(sn), consumable); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if consumable, ok := consumables[vacuum.ConsumableBrushSide]; ok {
			if e := b.MQTTPublishAsync(ctx, b.config.TopicConsumableBrushSide.Format(sn), consumable); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if consumable, ok := consumables[vacuum.ConsumableSensor]; ok {
			if e := b.MQTTPublishAsync(ctx, b.config.TopicConsumableSensor.Format(sn), consumable); e != nil {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	if e := b.updateState(ctx); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.updateFanPower(ctx); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.updateVolume(ctx); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}
