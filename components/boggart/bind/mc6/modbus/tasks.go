package modbus

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/multierr"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/mc6"
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
			WithName("read-write-registers-updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskReadWriteRegistersUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.StatusUpdaterInterval)),
	)

	if e != nil {
		err = multierr.Append(err, e)
	}

	_, e = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName("read-only-registers-updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskReadOnlyRegistersUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.SensorUpdaterInterval)),
	)

	if e != nil {
		err = multierr.Append(err, e)
	}

	return err
}

func (b *Bind) taskReadWriteRegistersUpdaterHandler(ctx context.Context) error {
	deviceType, err := b.DeviceType(ctx)
	if err != nil {
		return fmt.Errorf("get device type failed: %w", err)
	}

	provider := b.Provider()
	registers, e := provider.ReadAsMap(mc6.AddressStatus, mc6.AddressHoldingTime)
	if err != nil {
		return multierr.Append(err, e)
	}

	id := b.Meta().ID()
	cfg := b.config()

	if k := mc6.AddressStatus; deviceType.IsSupported(k) {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicStatusState.Format(id), registers[k].Bool()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressSystemMode; deviceType.IsSupported(k) {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicSystemModeState.Format(id), registers[k].Uint()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressFanSpeed; deviceType.IsSupported(k) {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicFanSpeedState.Format(id), registers[k].Uint()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressTargetTemperature; deviceType.IsSupported(k) {
		metricTargetTemperature.With("id", id).Set(registers[k].Temperature())

		if e = b.MQTT().PublishAsync(ctx, cfg.TopicTargetTemperatureState.Format(id), registers[k].Temperature()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressAway; deviceType.IsSupported(k) {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicAwayState.Format(id), registers[k].Bool()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressAwayTemperature; deviceType.IsSupported(k) {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicAwayTemperatureState.Format(id), registers[k].Temperature()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	//
	//
	// TODO:
	//
	//

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

func (b *Bind) taskReadOnlyRegistersUpdaterHandler(ctx context.Context) error {
	deviceType, err := b.DeviceType(ctx)
	if err != nil {
		return fmt.Errorf("get device type failed: %w", err)
	}

	provider := b.Provider()
	registers, e := provider.ReadAsMap(mc6.AddressRoomTemperature, mc6.AddressFanSpeedNumbers)
	if e != nil {
		return multierr.Append(err, e)
	}

	id := b.Meta().ID()
	cfg := b.config()

	if k := mc6.AddressRoomTemperature; deviceType.IsSupported(k) {
		metricRoomTemperature.With("id", id).Set(registers[k].Temperature())

		if e = b.MQTT().PublishAsync(ctx, cfg.TopicRoomTemperature.Format(id), registers[k].Temperature()); e != nil {
			err = multierr.Append(err, fmt.Errorf("get room temperature failed: %w", e))
		}
	}

	if k := mc6.AddressFloorTemperature; deviceType.IsSupported(k) {
		metricFloorTemperature.With("id", id).Set(registers[k].Temperature())

		if e = b.MQTT().PublishAsync(ctx, cfg.TopicFloorTemperature.Format(id), registers[k].Temperature()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressHumidity; deviceType.IsSupported(k) {
		metricHumidity.With("id", id).Set(float64(registers[k].Uint()))

		if e = b.MQTT().PublishAsync(ctx, cfg.TopicHumidity.Format(id), registers[k].Uint()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressHeatingValve; deviceType.IsSupported(k) {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicHeatingValve.Format(id), registers[k].Bool()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressCoolingValve; deviceType.IsSupported(k) {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicCoolingValve.Format(id), registers[k].Bool()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressHeatingOutput; deviceType.IsSupported(k) {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicHeatingOutputStatus.Format(id), registers[k].Bool()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressHoldingFunction; deviceType.IsSupported(k) {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicHoldingFunction.Format(id), registers[k].Bool()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressFloorOverheat; deviceType.IsSupported(k) {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicFloorOverheat.Format(id), registers[k].Bool()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	if k := mc6.AddressFanSpeedNumbers; deviceType.IsSupported(k) {
		if e = b.MQTT().PublishAsync(ctx, cfg.TopicFanSpeedNumbers.Format(id), registers[k].Uint()); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}
