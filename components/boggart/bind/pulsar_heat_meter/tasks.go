package pulsar_heat_meter

import (
	"context"
	"math"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(b.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(time.Minute)
	taskStateUpdater.SetName("bind-pulsar-heat-meter-state-updater-" + b.SerialNumber())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskStateUpdater(ctx context.Context) (interface{}, error) {
	if _, err := b.provider.Version(); err != nil {
		b.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	b.UpdateStatus(boggart.DeviceStatusOnline)
	serialNumber := b.SerialNumber()

	if currentVal, err := b.provider.TemperatureIn(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.temperatureIn))
		if current != prev {
			atomic.StoreUint64(&b.temperatureIn, math.Float64bits(current))

			b.MQTTPublishAsync(ctx, MQTTTopicTemperatureIn.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := b.provider.TemperatureOut(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.temperatureOut))
		if current != prev {
			atomic.StoreUint64(&b.temperatureOut, math.Float64bits(current))

			b.MQTTPublishAsync(ctx, MQTTTopicTemperatureOut.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := b.provider.TemperatureDelta(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.temperatureDelta))
		if current != prev {
			atomic.StoreUint64(&b.temperatureDelta, math.Float64bits(current))

			b.MQTTPublishAsync(ctx, MQTTTopicTemperatureDelta.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := b.provider.Energy(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.energy))
		if current != prev {
			atomic.StoreUint64(&b.energy, math.Float64bits(current))

			b.MQTTPublishAsync(ctx, MQTTTopicEnergy.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := b.provider.Consumption(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.consumption))
		if current != prev {
			atomic.StoreUint64(&b.consumption, math.Float64bits(current))

			b.MQTTPublishAsync(ctx, MQTTTopicConsumption.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := b.provider.Capacity(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.capacity))
		if current != prev {
			atomic.StoreUint64(&b.capacity, math.Float64bits(current))

			b.MQTTPublishAsync(ctx, MQTTTopicCapacity.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	if currentVal, err := b.provider.Power(); err == nil {
		current := float64(currentVal)
		prev := math.Float64frombits(atomic.LoadUint64(&b.power))
		if current != prev {
			atomic.StoreUint64(&b.power, math.Float64bits(current))

			b.MQTTPublishAsync(ctx, MQTTTopicPower.Format(serialNumber), 0, true, current)
		}
	} else {
		// TODO: log
	}

	// inputs
	if currentVal, err := b.provider.PulseInput1(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&b.input1)
		if current != prev {
			atomic.StoreUint64(&b.input1, current)

			b.MQTTPublishAsync(ctx, MQTTTopicInputPulses.Format(serialNumber, 1), 0, true, current)
			b.MQTTPublishAsync(ctx, MQTTTopicInputVolume.Format(serialNumber, 1), 0, true, b.inputVolume(current, b.config.Input1Offset))
		}
	} else {
		// TODO: log
	}

	if currentVal, err := b.provider.PulseInput2(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&b.input2)
		if current != prev {
			atomic.StoreUint64(&b.input2, current)

			b.MQTTPublishAsync(ctx, MQTTTopicInputPulses.Format(serialNumber, 2), 0, true, current)
			b.MQTTPublishAsync(ctx, MQTTTopicInputVolume.Format(serialNumber, 2), 0, true, b.inputVolume(current, b.config.Input2Offset))
		}
	} else {
		// TODO: log
	}

	if currentVal, err := b.provider.PulseInput3(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&b.input3)
		if current != prev {
			atomic.StoreUint64(&b.input3, current)

			b.MQTTPublishAsync(ctx, MQTTTopicInputPulses.Format(serialNumber, 3), 0, true, current)
			b.MQTTPublishAsync(ctx, MQTTTopicInputVolume.Format(serialNumber, 3), 0, true, b.inputVolume(current, b.config.Input3Offset))
		}
	} else {
		// TODO: log
	}

	if currentVal, err := b.provider.PulseInput4(); err == nil {
		current := uint64(currentVal)
		prev := atomic.LoadUint64(&b.input4)
		if current != prev {
			atomic.StoreUint64(&b.input4, current)

			b.MQTTPublishAsync(ctx, MQTTTopicInputPulses.Format(serialNumber, 4), 0, true, current)
			b.MQTTPublishAsync(ctx, MQTTTopicInputVolume.Format(serialNumber, 4), 0, true, b.inputVolume(current, b.config.Input4Offset))
		}
	} else {
		// TODO: log
	}

	return nil, nil
}
