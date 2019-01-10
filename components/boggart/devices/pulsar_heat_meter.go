package devices

import (
	"context"
	"encoding/hex"
	"math"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

const (
	PulsarHeadMeterMQTTTopicTemperatureIn    mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_in"
	PulsarHeadMeterMQTTTopicTemperatureOut   mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_out"
	PulsarHeadMeterMQTTTopicTemperatureDelta mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/temperature_delta"
	PulsarHeadMeterMQTTTopicEnergy           mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/energy"
	PulsarHeadMeterMQTTTopicConsumption      mqtt.Topic = boggart.ComponentName + "/meter/pulsar/+/consumption"
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

	boggart.DeviceBindBase
	boggart.DeviceBindSerialNumber
	boggart.DeviceBindMQTT

	provider *pulsar.HeatMeter
	interval time.Duration
}

func NewPulsarHeadMeter(provider *pulsar.HeatMeter, interval time.Duration) *PulsarHeadMeter {
	device := &PulsarHeadMeter{
		provider: provider,
		interval: interval,

		temperatureIn:    math.MaxUint64,
		temperatureOut:   math.MaxUint64,
		temperatureDelta: math.MaxUint64,
		energy:           math.MaxUint64,
		consumption:      math.MaxUint64,
	}
	device.Init()
	device.SetSerialNumber(hex.EncodeToString(provider.Address()))

	return device
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

func (d *PulsarHeadMeter) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(d.interval)
	taskStateUpdater.SetName("bind-pulsar-heat-meter-state-updater-" + d.SerialNumber())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (d *PulsarHeadMeter) taskStateUpdater(ctx context.Context) (interface{}, error) {
	if _, err := d.provider.Version(); err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)
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

		prev := math.Float64frombits(atomic.LoadUint64(&d.temperatureOut))
		if current != prev {
			atomic.StoreUint64(&d.temperatureOut, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeadMeterMQTTTopicTemperatureOut.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	if current, err := d.TemperatureDelta(ctx); err == nil {
		metricHeatMeterPulsarTemperatureDelta.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.temperatureDelta))
		if current != prev {
			atomic.StoreUint64(&d.temperatureDelta, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeadMeterMQTTTopicTemperatureDelta.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	if current, err := d.Energy(ctx); err == nil {
		metricHeatMeterPulsarEnergy.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.energy))
		if current != prev {
			atomic.StoreUint64(&d.energy, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeadMeterMQTTTopicEnergy.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	if current, err := d.Consumption(ctx); err == nil {
		metricHeatMeterPulsarConsumption.With("serial_number", serialNumber).Set(current)

		prev := math.Float64frombits(atomic.LoadUint64(&d.consumption))
		if current != prev {
			atomic.StoreUint64(&d.consumption, math.Float64bits(current))

			d.MQTTPublishAsync(ctx, PulsarHeadMeterMQTTTopicConsumption.Format(serialNumber), 0, true, current)
		}
	} else {
		return nil, err
	}

	return nil, nil
}

func (d *PulsarHeadMeter) MQTTTopics() []mqtt.Topic {
	sn := d.SerialNumber()

	return []mqtt.Topic{
		mqtt.Topic(PulsarHeadMeterMQTTTopicTemperatureIn.Format(sn)),
		mqtt.Topic(PulsarHeadMeterMQTTTopicTemperatureOut.Format(sn)),
		mqtt.Topic(PulsarHeadMeterMQTTTopicTemperatureDelta.Format(sn)),
		mqtt.Topic(PulsarHeadMeterMQTTTopicEnergy.Format(sn)),
		mqtt.Topic(PulsarHeadMeterMQTTTopicConsumption.Format(sn)),
	}
}
