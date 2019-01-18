package pulsar

import (
	"context"
	"math"
	"sync/atomic"

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

	if currentVal, err := b.provider.TemperatureIn(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.temperatureIn))
		if current != prev {
			atomic.StoreUint64(&b.temperatureIn, math.Float64bits(current))
			metricTemperatureIn.With("serial_number", sn).Set(current)

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicTemperatureIn.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if currentVal, err := b.provider.TemperatureOut(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.temperatureOut))
		if current != prev {
			atomic.StoreUint64(&b.temperatureOut, math.Float64bits(current))
			metricTemperatureOut.With("serial_number", sn).Set(current)

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicTemperatureOut.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if currentVal, err := b.provider.TemperatureDelta(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.temperatureDelta))
		if current != prev {
			atomic.StoreUint64(&b.temperatureDelta, math.Float64bits(current))
			metricTemperatureDelta.With("serial_number", sn).Set(current)

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicTemperatureDelta.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if currentVal, err := b.provider.Energy(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.energy))
		if current != prev {
			atomic.StoreUint64(&b.energy, math.Float64bits(current))
			metricEnergy.With("serial_number", sn).Set(current)

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicEnergy.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if currentVal, err := b.provider.Consumption(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.consumption))
		if current != prev {
			atomic.StoreUint64(&b.consumption, math.Float64bits(current))
			metricConsumption.With("serial_number", sn).Set(current)

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicConsumption.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if currentVal, err := b.provider.Capacity(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.capacity))
		if current != prev {
			atomic.StoreUint64(&b.capacity, math.Float64bits(current))
			metricCapacity.With("serial_number", sn).Set(current)

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicCapacity.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if currentVal, err := b.provider.Power(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.power))
		if current != prev {
			atomic.StoreUint64(&b.power, math.Float64bits(current))
			metricPower.With("serial_number", sn).Set(current)

			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicPower.Format(snMQTT), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	// inputs
	if currentVal, err := b.provider.PulseInput1(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&b.input1)
		if current != prev {
			atomic.StoreUint64(&b.input1, current)

			metricInputPulses.With("serial_number", sn).With("input", "1").Set(float64(current))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputPulses.Format(snMQTT, 1), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}

			volume := b.inputVolume(current, b.config.Input1Offset)
			metricInputVolume.With("serial_number", sn).With("input", "1").Set(volume)
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputVolume.Format(snMQTT, 1), 0, true, volume); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if currentVal, err := b.provider.PulseInput2(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&b.input2)
		if current != prev {
			atomic.StoreUint64(&b.input2, current)

			metricInputPulses.With("serial_number", sn).With("input", "2").Set(float64(current))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputPulses.Format(snMQTT, 2), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}

			volume := b.inputVolume(current, b.config.Input2Offset)
			metricInputVolume.With("serial_number", sn).With("input", "2").Set(volume)
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputVolume.Format(snMQTT, 2), 0, true, volume); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	if currentVal, err := b.provider.PulseInput3(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&b.input3)
		if current != prev {
			atomic.StoreUint64(&b.input3, current)

			metricInputPulses.With("serial_number", sn).With("input", "3").Set(float64(current))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputPulses.Format(snMQTT, 3), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}

			volume := b.inputVolume(current, b.config.Input3Offset)
			metricInputVolume.With("serial_number", sn).With("input", "3").Set(volume)
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputVolume.Format(snMQTT, 3), 0, true, volume); err != nil {
				result = multierr.Append(result, err)
			}

		}
	} else {
		result = multierr.Append(result, err)
	}

	if currentVal, err := b.provider.PulseInput4(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&b.input4)
		if current != prev {
			atomic.StoreUint64(&b.input4, current)

			metricInputPulses.With("serial_number", sn).With("input", "4").Set(float64(current))
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputPulses.Format(snMQTT, 4), 0, true, current); err != nil {
				result = multierr.Append(result, err)
			}

			volume := b.inputVolume(current, b.config.Input4Offset)
			metricInputVolume.With("serial_number", sn).With("input", "4").Set(volume)
			if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicInputVolume.Format(snMQTT, 4), 0, true, volume); err != nil {
				result = multierr.Append(result, err)
			}
		}
	} else {
		result = multierr.Append(result, err)
	}

	return nil, result
}
