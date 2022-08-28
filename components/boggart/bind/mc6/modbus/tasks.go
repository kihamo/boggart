package modbus

import (
	"context"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/components/mqtt"
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

func (b *Bind) taskDeviceTypeHandler(ctx context.Context) error {
	deviceType, err := b.DeviceType(ctx)
	if err != nil {
		return fmt.Errorf("get device type failed: %w", err)
	}

	cfg := b.config()

	// mqtt sensors
	subscribers := make([]mqtt.Subscriber, 0)
	id := b.Meta().ID()

	if deviceType.IsSupportedStatus() {
		subscribers = append(subscribers, mqtt.NewSubscriber(cfg.TopicStatus.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetStatus(message.Bool())
			}),
		))
	}

	if e := b.MQTT().Subscribe(subscribers...); e != nil {
		err = multierr.Append(err, e)
	}

	// tasks
	_, e := b.Workers().RegisterTask(
		tasks.NewTask().
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
			WithName("state-updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskStateUpdaterHandler),
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

func (b *Bind) taskStateUpdaterHandler(ctx context.Context) (err error) {
	// TODO: check b.deviceType

	provider := b.Provider()
	id := b.Meta().ID()
	cfg := b.config()

	if val, e := provider.Status(); e == nil {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicStatusState.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get status failed: %w", e))
	}

	if val, e := provider.HeatingValve(); e == nil {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicHeatingValve.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get heating valve status failed: %w", e))
	}

	if val, e := provider.CoolingValve(); e == nil {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicCoolingValve.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get cooling valve status failed: %w", e))
	}

	if val, e := provider.HeatingOutput(); e == nil {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicHeatingOutputStatus.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get heating output status failed: %w", e))
	}

	if val, e := provider.HoldingFunction(); e == nil {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicHoldingFunction.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get holding function failed: %w", e))
	}

	if val, e := provider.FloorOverheat(); e == nil {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicFloorOverheat.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get floor overheat failed: %w", e))
	}

	if val, e := provider.FanSpeedNumbers(); e == nil {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicFanSpeedNumbers.Format(id), val); e != nil {
			err = multierr.Append(err, e)
		}
	} else {
		err = multierr.Append(err, fmt.Errorf("get fan speed mode failed: %w", e))
	}

	return err
}

func (b *Bind) taskSensorUpdaterHandler(ctx context.Context) (err error) {
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
