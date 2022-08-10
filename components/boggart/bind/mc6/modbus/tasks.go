package modbus

import (
	"context"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("device-type").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskDeviceTypeHandler),
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

func (b *Bind) taskDeviceTypeHandler(ctx context.Context) (err error) {
	cfg := b.config()

	deviceType, err := b.Provider().DeviceType()
	if err != nil {
		return fmt.Errorf("get device type failed: %w", err)
	}

	b.deviceType.Set(uint64(deviceType))

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicDeviceType.Format(b.Meta().ID()), deviceType); e != nil {
		err = multierr.Append(err, e)
	}

	_, e := b.Workers().RegisterTask(
		tasks.NewTask().
			WithName("set-defaults-config").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskSetDefaultsConfigHandler),
				),
			).
			WithSchedule(
				tasks.ScheduleWithSuccessLimit(
					tasks.ScheduleWithDuration(tasks.ScheduleNow(), time.Second*30),
					1,
				),
			),
	)

	if e != nil {
		err = multierr.Append(err, e)
	}

	_, e = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName("status-updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskStatusUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.StatusUpdaterInterval)),
	)

	if e != nil {
		err = multierr.Append(err, e)
	}

	_, e = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName("sensor-updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskSensorUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.SensorUpdaterInterval)),
	)

	if e != nil {
		err = multierr.Append(err, e)
	}

	return err
}

func (b *Bind) taskSetDefaultsConfigHandler(ctx context.Context) (err error) {
	cfg := b.config()

	if e := b.AwayTemperature(ctx, cfg.DefaultsAwayTemperature); e != nil {
		err = multierr.Append(err, fmt.Errorf("set default away temperature failed: %w", e))
	}

	return err
}

func (b *Bind) taskStatusUpdaterHandler(ctx context.Context) (err error) {
	// TODO: check b.deviceType

	provider := b.Provider()
	id := b.Meta().ID()
	cfg := b.config()

	if val, e := provider.HeatingOutputStatus(); e == nil {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicHeatingOutputStatus.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get heating output status failed: %w", e))
	}

	return err
}

func (b *Bind) taskSensorUpdaterHandler(ctx context.Context) (err error) {
	// TODO: check b.deviceType

	provider := b.Provider()
	id := b.Meta().ID()
	cfg := b.config()

	if val, e := provider.RoomTemperature(); e == nil {
		metricRoomTemperature.With("id", id).Set(val)

		if e = b.MQTT().PublishAsync(ctx, cfg.TopicRoomTemperature.Format(id), val); e != nil {
			err = multierr.Append(err, fmt.Errorf("get room temperature failed: %w", e))
		}
	} else {
		err = multierr.Append(err, e)
	}

	if val, e := provider.FloorTemperature(); e == nil {
		metricFlourTemperature.With("id", id).Set(val)

		if e = b.MQTT().PublishAsync(ctx, cfg.TopicFloorTemperature.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get floor temperature failed: %w", e))
	}

	if val, e := provider.Humidity(); e == nil {
		metricHumidity.With("id", id).Set(float64(val))

		if e = b.MQTT().PublishAsync(ctx, cfg.TopicHumidity.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get humidity failed: %w", e))
	}

	return err
}
