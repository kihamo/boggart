package devices

import (
	"context"
	"encoding/hex"
	"math"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

const (
	PulsarHeadMeterMQTTTopicTemperatureIn    boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/pulsar/+/temperature_in"
	PulsarHeadMeterMQTTTopicTemperatureOut   boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/pulsar/+/temperature_out"
	PulsarHeadMeterMQTTTopicTemperatureDelta boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/pulsar/+/temperature_delta"
	PulsarHeadMeterMQTTTopicEnergy           boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/pulsar/+/energy"
	PulsarHeadMeterMQTTTopicConsumption      boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/pulsar/+/consumption"
)

var (
	metricHeatMeterPulsarTemperatureIn    = snitch.NewGauge(boggart.ComponentName+"_device_heat_meter_pulsar_temperature_in_celsius", "Pulsar temperature in")
	metricHeatMeterPulsarTemperatureOut   = snitch.NewGauge(boggart.ComponentName+"_device_heat_meter_pulsar_temperature_out_celsius", "Pulsar temperature out")
	metricHeatMeterPulsarTemperatureDelta = snitch.NewGauge(boggart.ComponentName+"_device_heat_meter_pulsar_temperature_delta_celsius", "Pulsar temperature delta")
	metricHeatMeterPulsarEnergy           = snitch.NewGauge(boggart.ComponentName+"_device_heat_meter_pulsar_energy_gigacolories", "Pulsar energy")
	metricHeatMeterPulsarConsumption      = snitch.NewGauge(boggart.ComponentName+"_device_heat_meter_pulsar_consumption_cubic_metres_per_hour", "Pulsar consumption")
)

type PulsarHeadMeter struct {
	temperatureIn    uint64
	temperatureOut   uint64
	temperatureDelta uint64
	energy           uint64
	consumption      uint64

	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	provider *pulsar.HeatMeter
	interval time.Duration
}

func NewPulsarHeadMeter(provider *pulsar.HeatMeter, interval time.Duration) *PulsarHeadMeter {
	device := &PulsarHeadMeter{
		provider: provider,
		interval: interval,
	}
	device.Init()
	device.SetSerialNumber(hex.EncodeToString(provider.Address()))
	device.SetDescription("Pulsar heat meter with serial number " + device.SerialNumber())

	return device
}

func (d *PulsarHeadMeter) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeHeatMeter,
	}
}

func (d *PulsarHeadMeter) TemperatureIn(context.Context) (float64, error) {
	value, err := d.provider.TemperatureIn()
	if err != nil {
		return -1, err
	}

	return float64(value), nil
}

func (d *PulsarHeadMeter) TemperatureOut(context.Context) (float64, error) {
	value, err := d.provider.TemperatureOut()
	if err != nil {
		return -1, err
	}

	return float64(value), nil
}

func (d *PulsarHeadMeter) TemperatureDelta(context.Context) (float64, error) {
	value, err := d.provider.TemperatureDelta()
	if err != nil {
		return -1, err
	}

	return float64(value), nil
}

func (d *PulsarHeadMeter) Energy(context.Context) (float64, error) {
	value, err := d.provider.Energy()
	if err != nil {
		return -1, err
	}

	return float64(value), nil
}

func (d *PulsarHeadMeter) Consumption(context.Context) (float64, error) {
	value, err := d.provider.Consumption()
	if err != nil {
		return -1, err
	}

	return float64(value), nil
}

func (d *PulsarHeadMeter) Describe(ch chan<- *snitch.Description) {
	serialNumber := d.SerialNumber()

	metricHeatMeterPulsarTemperatureIn.With("serial_number", serialNumber).Describe(ch)
	metricHeatMeterPulsarTemperatureOut.With("serial_number", serialNumber).Describe(ch)
	metricHeatMeterPulsarTemperatureDelta.With("serial_number", serialNumber).Describe(ch)
	metricHeatMeterPulsarEnergy.With("serial_number", serialNumber).Describe(ch)
	metricHeatMeterPulsarConsumption.With("serial_number", serialNumber).Describe(ch)
}

func (d *PulsarHeadMeter) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()

	metricHeatMeterPulsarTemperatureIn.With("serial_number", serialNumber).Collect(ch)
	metricHeatMeterPulsarTemperatureOut.With("serial_number", serialNumber).Collect(ch)
	metricHeatMeterPulsarTemperatureDelta.With("serial_number", serialNumber).Collect(ch)
	metricHeatMeterPulsarEnergy.With("serial_number", serialNumber).Collect(ch)
	metricHeatMeterPulsarConsumption.With("serial_number", serialNumber).Collect(ch)
}

func (d *PulsarHeadMeter) Ping(_ context.Context) bool {
	_, err := d.provider.Version()
	return err == nil
}

func (d *PulsarHeadMeter) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-heat-meter-pulsar-updater-" + d.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *PulsarHeadMeter) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	serialNumber := d.SerialNumber()

	if current, err := d.TemperatureIn(ctx); err == nil {
		metricHeatMeterPulsarTemperatureIn.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.temperatureIn))
		if current != prev {
			atomic.StoreUint64(&d.temperatureIn, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeadMeterMQTTTopicTemperatureIn.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	if current, err := d.TemperatureOut(ctx); err == nil {
		metricHeatMeterPulsarTemperatureOut.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.temperatureIn))
		if current != prev {
			atomic.StoreUint64(&d.temperatureIn, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeadMeterMQTTTopicTemperatureOut.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	if current, err := d.TemperatureDelta(ctx); err == nil {
		metricHeatMeterPulsarTemperatureDelta.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.temperatureIn))
		if current != prev {
			atomic.StoreUint64(&d.temperatureIn, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeadMeterMQTTTopicTemperatureDelta.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	if current, err := d.Energy(ctx); err == nil {
		metricHeatMeterPulsarEnergy.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.temperatureIn))
		if current != prev {
			atomic.StoreUint64(&d.temperatureIn, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeadMeterMQTTTopicEnergy.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	if current, err := d.Consumption(ctx); err == nil {
		metricHeatMeterPulsarConsumption.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.temperatureIn))
		if current != prev {
			atomic.StoreUint64(&d.temperatureIn, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeadMeterMQTTTopicConsumption.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	return nil, nil
}

func (d *PulsarHeadMeter) MQTTTopics() []boggart.DeviceMQTTTopic {
	return []boggart.DeviceMQTTTopic{
		PulsarHeadMeterMQTTTopicTemperatureIn,
		PulsarHeadMeterMQTTTopicTemperatureOut,
		PulsarHeadMeterMQTTTopicTemperatureDelta,
		PulsarHeadMeterMQTTTopicEnergy,
		PulsarHeadMeterMQTTTopicConsumption,
	}
}
