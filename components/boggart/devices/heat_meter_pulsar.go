package devices

import (
	"context"
	"encoding/hex"
	"reflect"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

var (
	metricHeatMeterPulsarTemperatureIn    = snitch.NewGauge(boggart.ComponentName+"_device_heat_meter_pulsar_temperature_in_celsius", "Pulsar temperature in")
	metricHeatMeterPulsarTemperatureOut   = snitch.NewGauge(boggart.ComponentName+"_device_heat_meter_pulsar_temperature_out_celsius", "Pulsar temperature out")
	metricHeatMeterPulsarTemperatureDelta = snitch.NewGauge(boggart.ComponentName+"_device_heat_meter_pulsar_temperature_delta_celsius", "Pulsar temperature delta")
	metricHeatMeterPulsarEnergy           = snitch.NewGauge(boggart.ComponentName+"_device_heat_meter_pulsar_energy_gigacolories", "Pulsar energy")
	metricHeatMeterPulsarConsumption      = snitch.NewGauge(boggart.ComponentName+"_device_heat_meter_pulsar_consumption_cubic_metres_per_hour", "Pulsar consumption")
)

type PulsarHeadMeterChange struct {
	TemperatureIn    float64
	TemperatureOut   float64
	TemperatureDelta float64
	Energy           float64
	Consumption      float64
}

type PulsarHeadMeter struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber

	provider *pulsar.HeatMeter
	interval time.Duration

	mutex      sync.Mutex
	lastValues PulsarHeadMeterChange
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
	currentValues := PulsarHeadMeterChange{}

	if value, err := d.TemperatureIn(ctx); err == nil {
		currentValues.TemperatureIn = value
		metricHeatMeterPulsarTemperatureIn.With("serial_number", serialNumber).Set(value)
	} else {
		return nil, err
	}

	if value, err := d.TemperatureOut(ctx); err == nil {
		currentValues.TemperatureOut = value
		metricHeatMeterPulsarTemperatureOut.With("serial_number", serialNumber).Set(value)
	} else {
		return nil, err
	}

	if value, err := d.TemperatureDelta(ctx); err == nil {
		currentValues.TemperatureDelta = value
		metricHeatMeterPulsarTemperatureDelta.With("serial_number", serialNumber).Set(value)
	} else {
		return nil, err
	}

	if value, err := d.Energy(ctx); err == nil {
		currentValues.Energy = value
		metricHeatMeterPulsarEnergy.With("serial_number", serialNumber).Set(value)
	} else {
		return nil, err
	}

	if value, err := d.Consumption(ctx); err == nil {
		currentValues.Consumption = value
		metricHeatMeterPulsarConsumption.With("serial_number", serialNumber).Set(value)
	} else {
		return nil, err
	}

	d.mutex.Lock()
	if !reflect.DeepEqual(d.lastValues, currentValues) {
		d.lastValues = currentValues
		d.mutex.Unlock()

		d.TriggerEvent(ctx, boggart.DeviceEventPulsarChanged, currentValues, serialNumber)
	} else {
		d.mutex.Unlock()
	}

	return nil, nil
}
