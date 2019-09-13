package pulsar

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater-" + b.SerialNumber())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) (interface{}, error) {
	if _, err := b.provider.Version(); err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	sn := b.SerialNumber()
	var result error

	if current, err := b.provider.TemperatureIn(); err == nil {
		metricTemperatureIn.With("serial_number", sn).Set(float64(current))

		if err := b.MQTTPublishAsync(ctx, b.config.TopicTemperatureIn, current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.TemperatureOut(); err == nil {
		metricTemperatureOut.With("serial_number", sn).Set(float64(current))

		if err := b.MQTTPublishAsync(ctx, b.config.TopicTemperatureOut, current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.TemperatureDelta(); err == nil {
		metricTemperatureDelta.With("serial_number", sn).Set(float64(current))

		if err := b.MQTTPublishAsync(ctx, b.config.TopicTemperatureDelta, current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Energy(); err == nil {
		metricEnergy.With("serial_number", sn).Set(float64(current))

		if err := b.MQTTPublishAsync(ctx, b.config.TopicEnergy, current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Consumption(); err == nil {
		metricConsumption.With("serial_number", sn).Set(float64(current))

		if err := b.MQTTPublishAsync(ctx, b.config.TopicConsumption, current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Capacity(); err == nil {
		metricCapacity.With("serial_number", sn).Set(float64(current))

		if err := b.MQTTPublishAsync(ctx, b.config.TopicCapacity, current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Power(); err == nil {
		metricPower.With("serial_number", sn).Set(float64(current))

		if err := b.MQTTPublishAsync(ctx, b.config.TopicPower, current); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	// inputs
	if current, err := b.provider.PulseInput1(); err == nil {
		metricInputPulses.With("serial_number", sn).With("input", "1").Set(float64(current))
		if err := b.MQTTPublishAsync(ctx, b.config.TopicInputPulses1, current); err != nil {
			result = multierr.Append(result, err)
		}

		volume := b.inputVolume(current, b.config.Input1Offset)
		metricInputVolume.With("serial_number", sn).With("input", "1").Set(float64(volume))
		if err := b.MQTTPublishAsync(ctx, b.config.TopicInputVolume1, volume); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.PulseInput2(); err == nil {
		metricInputPulses.With("serial_number", sn).With("input", "2").Set(float64(current))
		if err := b.MQTTPublishAsync(ctx, b.config.TopicInputPulses2, current); err != nil {
			result = multierr.Append(result, err)
		}

		volume := b.inputVolume(current, b.config.Input2Offset)
		metricInputVolume.With("serial_number", sn).With("input", "2").Set(float64(volume))
		if err := b.MQTTPublishAsync(ctx, b.config.TopicInputVolume2, volume); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.PulseInput3(); err == nil {
		metricInputPulses.With("serial_number", sn).With("input", "3").Set(float64(current))
		if err := b.MQTTPublishAsync(ctx, b.config.TopicInputPulses3, current); err != nil {
			result = multierr.Append(result, err)
		}

		volume := b.inputVolume(current, b.config.Input3Offset)
		metricInputVolume.With("serial_number", sn).With("input", "3").Set(float64(volume))
		if err := b.MQTTPublishAsync(ctx, b.config.TopicInputVolume3, volume); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.PulseInput4(); err == nil {
		metricInputPulses.With("serial_number", sn).With("input", "4").Set(float64(current))
		if err := b.MQTTPublishAsync(ctx, b.config.TopicInputPulses4, current); err != nil {
			result = multierr.Append(result, err)
		}

		volume := b.inputVolume(current, b.config.Input4Offset)
		metricInputVolume.With("serial_number", sn).With("input", "4").Set(float64(volume))
		if err := b.MQTTPublishAsync(ctx, b.config.TopicInputVolume4, volume); err != nil {
			result = multierr.Append(result, err)
		}
	} else {
		result = multierr.Append(result, err)
	}

	return nil, result
}
