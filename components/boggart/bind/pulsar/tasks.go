package pulsar

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.updaterInterval)
	taskStateUpdater.SetName("bind-pulsar-heat-meter-updater-" + b.SerialNumber())

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
	snMQTT := mqtt.NameReplace(sn)
	var result error

	if current, err := b.provider.TemperatureIn(); err == nil {
		if ok := b.temperatureIn.Set(current); ok {
			metricTemperatureIn.With("serial_number", sn).Set(float64(current))

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicTemperatureIn.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.TemperatureOut(); err == nil {
		if ok := b.temperatureOut.Set(current); ok {
			metricTemperatureOut.With("serial_number", sn).Set(float64(current))

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicTemperatureOut.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.TemperatureDelta(); err == nil {
		if ok := b.temperatureDelta.Set(current); ok {
			metricTemperatureDelta.With("serial_number", sn).Set(float64(current))

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicTemperatureDelta.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Energy(); err == nil {
		if ok := b.energy.Set(current); ok {
			metricEnergy.With("serial_number", sn).Set(float64(current))

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicEnergy.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Consumption(); err == nil {
		if ok := b.consumption.Set(current); ok {
			metricConsumption.With("serial_number", sn).Set(float64(current))

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicConsumption.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Capacity(); err == nil {
		if ok := b.capacity.Set(current); ok {
			metricCapacity.With("serial_number", sn).Set(float64(current))

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicCapacity.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.Power(); err == nil {
		if ok := b.power.Set(current); ok {
			metricPower.With("serial_number", sn).Set(float64(current))

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicPower.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	// inputs
	if current, err := b.provider.PulseInput1(); err == nil {
		if ok := b.input1.Set(current); ok {
			metricInputPulses.With("serial_number", sn).With("input", "1").Set(float64(current))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputPulses.Format(snMQTT, 1), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}

			volume := b.inputVolume(current, b.config.Input1Offset)
			metricInputVolume.With("serial_number", sn).With("input", "1").Set(float64(volume))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputVolume.Format(snMQTT, 1), 0, true, volume); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.PulseInput2(); err == nil {
		if ok := b.input2.Set(current); ok {
			metricInputPulses.With("serial_number", sn).With("input", "2").Set(float64(current))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputPulses.Format(snMQTT, 2), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}

			volume := b.inputVolume(current, b.config.Input2Offset)
			metricInputVolume.With("serial_number", sn).With("input", "2").Set(float64(volume))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputVolume.Format(snMQTT, 2), 0, true, volume); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.PulseInput3(); err == nil {
		if ok := b.input3.Set(current); ok {
			metricInputPulses.With("serial_number", sn).With("input", "3").Set(float64(current))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputPulses.Format(snMQTT, 3), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}

			volume := b.inputVolume(current, b.config.Input3Offset)
			metricInputVolume.With("serial_number", sn).With("input", "3").Set(float64(volume))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputVolume.Format(snMQTT, 3), 0, true, volume); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if current, err := b.provider.PulseInput4(); err == nil {
		if ok := b.input4.Set(current); ok {
			metricInputPulses.With("serial_number", sn).With("input", "4").Set(float64(current))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputPulses.Format(snMQTT, 4), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}

			volume := b.inputVolume(current, b.config.Input4Offset)
			metricInputVolume.With("serial_number", sn).With("input", "4").Set(float64(volume))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputVolume.Format(snMQTT, 4), 0, true, volume); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	return nil, result
}
