package pulsar

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/pulsar"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	list := make([]tasks.Task, 0, 2)

	if b.config.Address == "" {
		list = append(list,
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
		)
	} else {
		list = append(list, tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.UpdaterInterval)),
		)
	}

	return list
}

func (b *Bind) taskSerialNumberHandler(ctx context.Context) error {
	conn, err := b.getConnection()
	if err != nil {
		return err
	}

	address, err := pulsar.DeviceAddress(conn)
	if err != nil {
		return err
	}

	if err = b.createProvider(address); err != nil {
		return err
	}

	_, err = b.Workers().RegisterTask(tasks.NewTask().
		WithName("updater").
		WithHandler(
			b.Workers().WrapTaskHandlerIsOnline(
				tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
			),
		).
		WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config.UpdaterInterval)),
	)

	return err
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return nil
	}

	provider := b.Provider()
	if provider == nil {
		return nil
	}

	var result error

	if current, err := provider.TemperatureIn(); err == nil {
		metricTemperatureIn.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicTemperatureIn.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := provider.TemperatureOut(); err == nil {
		metricTemperatureOut.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicTemperatureOut.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := provider.TemperatureDelta(); err == nil {
		metricTemperatureDelta.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicTemperatureDelta.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := provider.Energy(); err == nil {
		metricEnergy.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicEnergy.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := provider.Consumption(); err == nil {
		metricConsumption.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicConsumption.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := provider.Capacity(); err == nil {
		metricCapacity.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicCapacity.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := provider.Power(); err == nil {
		metricPower.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicPower.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	// inputs
	if b.config.InputsCount > 0 {
		if current, err := provider.PulseInput1(); err == nil {
			const input1 = "1"

			volume := b.inputVolume(current, b.config.Input1Offset)

			metricInputPulses.With("serial_number", sn).With("input", input1).Set(float64(current))
			metricInputVolume.With("serial_number", sn).With("input", input1).Set(float64(volume))

			if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputPulses.Format(sn, input1), current); err != nil {
				result = multierr.Append(result, err)
			}

			if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputVolume.Format(sn, input1), volume); err != nil {
				result = multierr.Append(result, err)
			}
		} else {
			result = multierr.Append(result, err)
		}
	}

	if b.config.InputsCount > 1 {
		if current, err := provider.PulseInput2(); err == nil {
			const input2 = "2"

			volume := b.inputVolume(current, b.config.Input2Offset)

			metricInputPulses.With("serial_number", sn).With("input", input2).Set(float64(current))
			metricInputVolume.With("serial_number", sn).With("input", input2).Set(float64(volume))

			if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputPulses.Format(sn, input2), current); err != nil {
				result = multierr.Append(result, err)
			}

			if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputVolume.Format(sn, input2), volume); err != nil {
				result = multierr.Append(result, err)
			}
		} else {
			result = multierr.Append(result, err)
		}
	}

	if b.config.InputsCount > 2 {
		if current, err := provider.PulseInput3(); err == nil {
			const input3 = "3"

			volume := b.inputVolume(current, b.config.Input3Offset)

			metricInputPulses.With("serial_number", sn).With("input", input3).Set(float64(current))
			metricInputVolume.With("serial_number", sn).With("input", input3).Set(float64(volume))

			if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputPulses.Format(sn, input3), current); err != nil {
				result = multierr.Append(result, err)
			}

			if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputVolume.Format(sn, input3), volume); err != nil {
				result = multierr.Append(result, err)
			}
		} else {
			result = multierr.Append(result, err)
		}
	}

	if b.config.InputsCount > 3 {
		if current, err := provider.PulseInput4(); err == nil {
			const input4 = "4"

			volume := b.inputVolume(current, b.config.Input4Offset)

			metricInputPulses.With("serial_number", sn).With("input", input4).Set(float64(current))
			metricInputVolume.With("serial_number", sn).With("input", input4).Set(float64(volume))

			if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputPulses.Format(sn, input4), current); err != nil {
				result = multierr.Append(result, err)
			}

			if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputVolume.Format(sn, input4), volume); err != nil {
				result = multierr.Append(result, err)
			}
		} else {
			result = multierr.Append(result, err)
		}
	}

	return result
}
