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

	if deviceType.IsSupportedSystemMode() {
		subscribers = append(subscribers, mqtt.NewSubscriber(cfg.TopicSystemMode.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetSystemMode(uint16(message.Uint64()))
			}),
		))
	}

	if deviceType.IsSupportedFanSpeed() {
		subscribers = append(subscribers, mqtt.NewSubscriber(cfg.TopicFanSpeed.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetFanSpeed(uint16(message.Uint64()))
			}),
		))
	}

	if deviceType.IsSupportedTargetTemperature() {
		subscribers = append(subscribers, mqtt.NewSubscriber(cfg.TopicTargetTemperature.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetTargetTemperature(message.Float64())
			}),
		))
	}

	if deviceType.IsSupportedAway() {
		subscribers = append(subscribers, mqtt.NewSubscriber(cfg.TopicAway.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetAway(message.Bool())
			}),
		))
	}

	if deviceType.IsSupportedAwayTemperature() {
		subscribers = append(subscribers, mqtt.NewSubscriber(cfg.TopicAwayTemperature.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetAwayTemperature(uint16(message.Uint64()))
			}),
		))
	}

	if deviceType.IsSupportedHoldingTemperatureAndTime() {
		subscribers = append(subscribers, mqtt.NewSubscriber(cfg.TopicHoldingTime.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetHoldingTime(time.Duration(message.Uint64()) * time.Minute)
			}),
		))
	}

	if deviceType.IsSupportedHoldingTemperature() {
		subscribers = append(subscribers, mqtt.NewSubscriber(cfg.TopicHoldingTemperature.Format(id), 0,
			b.MQTTWrapSubscriberChangeState(func(message mqtt.Message) error {
				return b.Provider().SetHoldingTemperature(uint16(message.Uint64()))
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

func (b *Bind) taskStateUpdaterHandler(ctx context.Context) error {
	deviceType, err := b.DeviceType(ctx)
	if err != nil {
		return fmt.Errorf("get device type failed: %w", err)
	}

	provider := b.Provider()
	id := b.Meta().ID()
	cfg := b.config()

	if deviceType.IsSupportedStatus() {
		if val, e := provider.Status(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicStatusState.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get status failed: %w", e))
		}
	}

	if deviceType.IsSupportedHeatingValve() {
		if val, e := provider.HeatingValve(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicHeatingValve.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get heating valve status failed: %w", e))
		}
	}

	if deviceType.IsSupportedCoolingValve() {
		if val, e := provider.CoolingValve(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicCoolingValve.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get cooling valve status failed: %w", e))
		}
	}

	if deviceType.IsSupportedHeatingOutput() {
		if val, e := provider.HeatingOutput(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicHeatingOutputStatus.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get heating output status failed: %w", e))
		}
	}

	if deviceType.IsSupportedHoldingFunction() {
		if val, e := provider.HoldingFunction(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicHoldingFunction.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get holding function failed: %w", e))
		}
	}

	if deviceType.IsSupportedFloorOverheat() {
		if val, e := provider.FloorOverheat(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicFloorOverheat.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get floor overheat failed: %w", e))
		}
	}

	if deviceType.IsSupportedFanSpeedNumbers() {
		if val, e := provider.FanSpeedNumbers(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicFanSpeedNumbers.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get fan speed mode failed: %w", e))
		}
	}

	if deviceType.IsSupportedSystemMode() {
		if val, e := provider.SystemMode(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicSystemModeState.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get system mode state failed: %w", e))
		}
	}

	if deviceType.IsSupportedFanSpeed() {
		if val, e := provider.FanSpeed(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicFanSpeedState.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get fan speed state failed: %w", e))
		}
	}

	if deviceType.IsSupportedTargetTemperature() {
		if val, e := provider.TargetTemperature(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicTargetTemperatureState.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get target temperature state failed: %w", e))
		}
	}

	if deviceType.IsSupportedAway() {
		if val, e := provider.Away(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicAwayState.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get away state failed: %w", e))
		}
	}

	if deviceType.IsSupportedAwayTemperature() {
		if val, e := provider.AwayTemperature(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicAwayTemperatureState.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get away temperature state failed: %w", e))
		}
	}

	if deviceType.IsSupportedHoldingTemperatureAndTime() {
		if temp, tm, e := provider.HoldingTemperatureAndTime(); e == nil {
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicHoldingTemperatureState.Format(id), temp); e != nil {
				err = multierr.Append(err, e)
			}
			if e = b.MQTT().PublishAsync(ctx, cfg.TopicHoldingTimeState.Format(id), tm.Minutes()); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get holding temperature and time state failed: %w", e))
		}
	} else {
		if deviceType.IsSupportedHoldingTime() {
			if val, e := provider.HoldingTime(); e == nil {
				if e = b.MQTT().PublishAsync(ctx, cfg.TopicHoldingTimeState.Format(id), val.Minutes()); e != nil {
					err = multierr.Append(err, e)
				}
			} else {
				err = multierr.Append(err, fmt.Errorf("get holding time state failed: %w", e))
			}
		}

		if deviceType.IsSupportedHoldingTemperature() {
			if val, e := provider.HoldingTemperature(); e == nil {
				if e = b.MQTT().PublishAsync(ctx, cfg.TopicHoldingTemperatureState.Format(id), val); e != nil {
					err = multierr.Append(err, e)
				}
			} else {
				err = multierr.Append(err, fmt.Errorf("get holding temperature state failed: %w", e))
			}
		}
	}

	return err
}

func (b *Bind) taskSensorUpdaterHandler(ctx context.Context) error {
	deviceType, err := b.DeviceType(ctx)
	if err != nil {
		return fmt.Errorf("get device type failed: %w", err)
	}

	provider := b.Provider()
	id := b.Meta().ID()
	cfg := b.config()

	if deviceType.IsSupportedRoomTemperature() {
		if val, e := provider.RoomTemperature(); e == nil {
			metricRoomTemperature.With("id", id).Set(val)

			if e = b.MQTT().PublishAsync(ctx, cfg.TopicRoomTemperature.Format(id), val); e != nil {
				err = multierr.Append(err, fmt.Errorf("get room temperature failed: %w", e))
			}
		} else {
			err = multierr.Append(err, e)
		}
	}

	if deviceType.IsSupportedFloorTemperature() {
		if val, e := provider.FloorTemperature(); e == nil {
			metricFloorTemperature.With("id", id).Set(val)

			if e = b.MQTT().PublishAsync(ctx, cfg.TopicFloorTemperature.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get floor temperature failed: %w", e))
		}
	}

	if deviceType.IsSupportedHumidity() {
		if val, e := provider.Humidity(); e == nil {
			metricHumidity.With("id", id).Set(float64(val))

			if e = b.MQTT().PublishAsync(ctx, cfg.TopicHumidity.Format(id), val); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, fmt.Errorf("get humidity failed: %w", e))
		}
	}

	return err
}
