package miio

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/xiaomi/miio/devices/vacuum"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("serial-number").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskSerialNumberHandler),
				),
			).
			WithSchedule(
				tasks.ScheduleWithSuccessLimit(
					tasks.ScheduleWithDuration(tasks.ScheduleNow(), time.Second*30),
					1,
				),
			),
	}
}

func (b *Bind) taskSerialNumberHandler(ctx context.Context) error {
	sn, err := b.device.SerialNumber(ctx)
	if err != nil {
		return err
	}

	b.Meta().SetSerialNumber(sn)
	cfg := b.config()

	_, err = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerWithTimeout(
						tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
						cfg.UpdaterTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.UpdaterInterval)),
	)

	return err
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
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

	cfg := b.config()

	// only statistics
	status, e := b.device.Status(ctx)
	if e == nil {
		metricBattery.With("serial_number", sn).Set(float64(status.Battery))

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicBattery.Format(sn), status.Battery); e != nil {
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
				metricCleanArea.With("serial_number", sn).Set(float64(lastClean.Area))
				metricCleanTime.With("serial_number", sn).Set(lastClean.CleaningDuration.Duration.Seconds())

				if e := b.MQTT().PublishAsync(ctx, cfg.TopicLastCleanCompleted.Format(sn), lastClean.Completed); e != nil {
					err = multierr.Append(err, e)
				}

				if e := b.MQTT().PublishAsync(ctx, cfg.TopicLastCleanArea.Format(sn), uint64(lastClean.Area)); e != nil {
					err = multierr.Append(err, e)
				}

				if e := b.MQTT().PublishAsync(ctx, cfg.TopicLastCleanStartDateTime.Format(sn), lastClean.StartTime); e != nil {
					err = multierr.Append(err, e)
				}

				if e := b.MQTT().PublishAsync(ctx, cfg.TopicLastCleanEndDateTime.Format(sn), lastClean.EndTime); e != nil {
					err = multierr.Append(err, e)
				}

				if e := b.MQTT().PublishAsync(ctx, cfg.TopicLastCleanDuration.Format(sn), lastClean.CleaningDuration.Duration); e != nil {
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
			if e := b.MQTT().PublishAsync(ctx, cfg.TopicConsumableFilter.Format(sn), consumable); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if consumable, ok := consumables[vacuum.ConsumableBrushMain]; ok {
			if e := b.MQTT().PublishAsync(ctx, cfg.TopicConsumableBrushMain.Format(sn), consumable); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if consumable, ok := consumables[vacuum.ConsumableBrushSide]; ok {
			if e := b.MQTT().PublishAsync(ctx, cfg.TopicConsumableBrushSide.Format(sn), consumable); e != nil {
				err = multierr.Append(err, e)
			}
		}

		if consumable, ok := consumables[vacuum.ConsumableSensor]; ok {
			if e := b.MQTT().PublishAsync(ctx, cfg.TopicConsumableSensor.Format(sn), consumable); e != nil {
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
