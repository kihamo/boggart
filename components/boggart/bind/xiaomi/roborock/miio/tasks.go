package miio

import (
	"context"

	"github.com/kihamo/boggart/providers/xiaomi/miio/devices/vacuum"
	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskState := b.Workers().WrapTaskIsOnline(b.taskUpdater)
	taskState.SetTimeout(b.config.UpdaterTimeout)
	taskState.SetRepeats(-1)
	taskState.SetRepeatInterval(b.config.UpdaterInterval)
	taskState.SetName("updater")

	return []workers.Task{
		taskState,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return nil
	}

	var err error

	if b.Meta().MAC() == nil {
		info, e := b.device.Info(ctx)
		if e != nil {
			return e
		}

		b.Meta().SetMAC(info.MAC.HardwareAddr)
	}

	// only statistics
	status, e := b.device.Status(ctx)
	if e == nil {
		metricBattery.With("serial_number", sn).Set(float64(status.Battery))

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicBattery.Format(sn), status.Battery); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, e)
	}

	summary, e := b.device.CleanSummary(ctx)
	if e == nil {
		if len(summary.CleanupIDs) > 0 {
			lastClean, e := b.device.CleanDetails(ctx, summary.CleanupIDs[0])
			if e == nil {
				if e := b.MQTT().PublishAsync(ctx, b.config.TopicLastCleanCompleted.Format(sn), lastClean.Completed); e != nil {
					err = multierr.Append(err, e)
				}

				if e := b.MQTT().PublishAsync(ctx, b.config.TopicLastCleanArea.Format(sn), lastClean.Area); e != nil {
					err = multierr.Append(err, e)
				}

				if e := b.MQTT().PublishAsync(ctx, b.config.TopicLastCleanStartDateTime.Format(sn), lastClean.StartTime); e != nil {
					err = multierr.Append(err, e)
				}

				if e := b.MQTT().PublishAsync(ctx, b.config.TopicLastCleanStartEndTime.Format(sn), lastClean.EndTime); e != nil {
					err = multierr.Append(err, e)
				}

				if e := b.MQTT().PublishAsync(ctx, b.config.TopicLastCleanDuration.Format(sn), lastClean.CleaningDuration); e != nil {
					err = multierr.Append(err, e)
				}
			} else {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	consumables, e := b.device.Consumables(ctx)
	if e == nil {
		if consumable, ok := consumables[vacuum.ConsumableFilter]; ok {
			if e := b.MQTT().PublishAsync(ctx, b.config.TopicConsumableFilter.Format(sn), consumable); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if consumable, ok := consumables[vacuum.ConsumableBrushMain]; ok {
			if e := b.MQTT().PublishAsync(ctx, b.config.TopicConsumableBrushMain.Format(sn), consumable); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if consumable, ok := consumables[vacuum.ConsumableBrushSide]; ok {
			if e := b.MQTT().PublishAsync(ctx, b.config.TopicConsumableBrushSide.Format(sn), consumable); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if consumable, ok := consumables[vacuum.ConsumableSensor]; ok {
			if e := b.MQTT().PublishAsync(ctx, b.config.TopicConsumableSensor.Format(sn), consumable); e != nil {
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

	// last clean time
	if summary, e := b.device.CleanSummary(ctx); e == nil {
		if len(summary.CleanupIDs) > 0 {
			if details, e := b.device.CleanDetails(ctx, summary.CleanupIDs[0]); e == nil {
				if e := b.MQTT().PublishAsync(ctx, b.config.TopicLastCleanStartDateTime.Format(sn), details.EndTime); e != nil {
					err = multierr.Append(err, e)
				}
			} else {
				err = multierr.Append(err, e)
			}
		}
	} else {
		err = multierr.Append(err, e)
	}

	return err
}
