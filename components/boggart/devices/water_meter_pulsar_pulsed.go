package devices

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

const (
	PulsarPulsedWaterMeterScale = 1000
)

var (
	metricWaterMeterPulsarPulsedVolume = snitch.NewGauge(boggart.ComponentName+"_device_water_meter_pulsar_pulsed_volume_cubic_metres", "Pulsar volume of water in cubic metres")
	metricWaterMeterPulsarPulsedPulses = snitch.NewGauge(boggart.ComponentName+"_device_water_meter_pulsar_pulsed_pulses", "Pulsar volume of water in pulses")
)

type PulsarPulsedWaterMeterChanged struct {
	Volume float64
	Pulses uint64
}

type PulsarPulsedWaterMeter struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber

	input        uint64
	volumeOffset float64
	provider     *pulsar.HeatMeter
	interval     time.Duration

	mutex      sync.Mutex
	lastValues PulsarPulsedWaterMeterChanged
}

func NewPulsarPulsedWaterMeter(serialNumber string, volumeOffset float64, provider *pulsar.HeatMeter, input uint64, interval time.Duration) *PulsarPulsedWaterMeter {
	device := &PulsarPulsedWaterMeter{
		volumeOffset: volumeOffset,
		provider:     provider,
		input:        input,
		interval:     interval,
	}
	device.Init()
	device.SetSerialNumber(serialNumber)

	return device
}

func (d *PulsarPulsedWaterMeter) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeWaterMeter,
	}
}

func (d *PulsarPulsedWaterMeter) Volume(ctx context.Context) (float64, error) {
	pulses, err := d.Pulses(ctx)
	if err != nil {
		return -1, err
	}

	return d.volume(pulses), nil
}

func (d *PulsarPulsedWaterMeter) Pulses(_ context.Context) (uint64, error) {
	var getFunc func() (float32, error)

	switch d.input {
	case pulsar.Input1:
		getFunc = d.provider.PulseInput1
	case pulsar.Input2:
		getFunc = d.provider.PulseInput2
	}

	pulses, err := getFunc()
	if err != nil {
		return 0, err
	}

	return uint64(pulses), nil
}

func (d *PulsarPulsedWaterMeter) Describe(ch chan<- *snitch.Description) {
	serialNumber := d.SerialNumber()

	metricWaterMeterPulsarPulsedVolume.With("serial_number", serialNumber).Describe(ch)
	metricWaterMeterPulsarPulsedPulses.With("serial_number", serialNumber).Describe(ch)
}

func (d *PulsarPulsedWaterMeter) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()

	metricWaterMeterPulsarPulsedVolume.With("serial_number", serialNumber).Collect(ch)
	metricWaterMeterPulsarPulsedPulses.With("serial_number", serialNumber).Collect(ch)
}

func (d *PulsarPulsedWaterMeter) Ping(_ context.Context) bool {
	_, err := d.provider.Version()
	return err == nil
}

func (d *PulsarPulsedWaterMeter) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-water-meter-pulsar-pulsed-updater-" + d.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *PulsarPulsedWaterMeter) volume(pulses uint64) float64 {
	return (d.volumeOffset*PulsarPulsedWaterMeterScale + float64(pulses*10)) / PulsarPulsedWaterMeterScale
}

func (d *PulsarPulsedWaterMeter) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	pulses, err := d.Pulses(ctx)
	if err != nil {
		return nil, err
	}

	serialNumber := d.SerialNumber()
	currentValues := PulsarPulsedWaterMeterChanged{
		Volume: d.volume(pulses),
		Pulses: pulses,
	}

	metricWaterMeterPulsarPulsedPulses.With("serial_number", serialNumber).Set(float64(currentValues.Pulses))
	metricWaterMeterPulsarPulsedVolume.With("serial_number", serialNumber).Set(currentValues.Volume)

	d.mutex.Lock()
	if !reflect.DeepEqual(d.lastValues, currentValues) {
		d.lastValues = currentValues
		d.mutex.Unlock()

		d.TriggerEvent(ctx, boggart.DeviceEventPulsarPulsedChanged, currentValues, serialNumber)
	} else {
		d.mutex.Unlock()
	}

	return nil, nil
}
