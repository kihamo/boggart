package pulsar

import (
	"context"
	"time"

	"github.com/kihamo/boggart/providers/pulsar"
	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	tasks := make([]workers.Task, 0, 2)

	if b.config.Address == "" {
		taskSerialNumber := b.Workers().WrapTaskOnceSuccess(b.taskSerialNumber)
		taskSerialNumber.SetRepeats(-1)
		taskSerialNumber.SetRepeatInterval(time.Second * 30)
		taskSerialNumber.SetName("serial-number")
		tasks = append(tasks, taskSerialNumber)
	}

	taskStateUpdater := b.Workers().WrapTaskIsOnline(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater")
	tasks = append(tasks, taskStateUpdater)

	return tasks
}

func (b *Bind) taskSerialNumber(ctx context.Context) error {
	address, err := pulsar.DeviceAddress(b.connection)
	if err != nil {
		return err
	}

	b.createProvider(address)
	return nil
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return nil
	}

	var result error

	if current, err := b.provider.TemperatureIn(); err == nil {
		metricTemperatureIn.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicTemperatureIn.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.TemperatureOut(); err == nil {
		metricTemperatureOut.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicTemperatureOut.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.TemperatureDelta(); err == nil {
		metricTemperatureDelta.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicTemperatureDelta.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Energy(); err == nil {
		metricEnergy.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicEnergy.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Consumption(); err == nil {
		metricConsumption.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicConsumption.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Capacity(); err == nil {
		metricCapacity.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicCapacity.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Power(); err == nil {
		metricPower.With("serial_number", sn).Set(float64(current))

		if err := b.MQTT().PublishAsync(ctx, b.config.TopicPower.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	// inputs
	if current, err := b.provider.PulseInput1(); err == nil {
		metricInputPulses.With("serial_number", sn).With("input", "1").Set(float64(current))
		if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputPulses1.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}

		volume := b.inputVolume(current, b.config.Input1Offset)
		metricInputVolume.With("serial_number", sn).With("input", "1").Set(float64(volume))
		if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputVolume1.Format(sn), volume); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.PulseInput2(); err == nil {
		metricInputPulses.With("serial_number", sn).With("input", "2").Set(float64(current))
		if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputPulses2.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}

		volume := b.inputVolume(current, b.config.Input2Offset)
		metricInputVolume.With("serial_number", sn).With("input", "2").Set(float64(volume))
		if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputVolume2.Format(sn), volume); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.PulseInput3(); err == nil {
		metricInputPulses.With("serial_number", sn).With("input", "3").Set(float64(current))
		if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputPulses3.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}

		volume := b.inputVolume(current, b.config.Input3Offset)
		metricInputVolume.With("serial_number", sn).With("input", "3").Set(float64(volume))
		if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputVolume3.Format(sn), volume); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.PulseInput4(); err == nil {
		metricInputPulses.With("serial_number", sn).With("input", "4").Set(float64(current))
		if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputPulses4.Format(sn), current); err != nil {
			result = multierr.Append(result, err)
		}

		volume := b.inputVolume(current, b.config.Input4Offset)
		metricInputVolume.With("serial_number", sn).With("input", "4").Set(float64(volume))
		if err := b.MQTT().PublishAsync(ctx, b.config.TopicInputVolume4.Format(sn), volume); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	return result
}
