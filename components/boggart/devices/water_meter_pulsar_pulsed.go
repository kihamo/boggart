package devices

import (
	"context"
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

type PulsarPulsedWaterMeter struct {
	boggart.DeviceBase

	input        uint64
	volumeOffset float64
	serialNumber string
	provider     *pulsar.HeatMeter
	interval     time.Duration
}

func NewPulsarPulsedWaterMeter(serialNumber string, volumeOffset float64, provider *pulsar.HeatMeter, input uint64, interval time.Duration) *PulsarPulsedWaterMeter {
	device := &PulsarPulsedWaterMeter{
		serialNumber: serialNumber,
		volumeOffset: volumeOffset,
		provider:     provider,
		input:        input,
		interval:     interval,
	}
	device.Init()

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
	metricWaterMeterPulsarPulsedVolume.With("serial_number", d.serialNumber).Describe(ch)
	metricWaterMeterPulsarPulsedPulses.With("serial_number", d.serialNumber).Describe(ch)
}

func (d *PulsarPulsedWaterMeter) Collect(ch chan<- snitch.Metric) {
	metricWaterMeterPulsarPulsedVolume.With("serial_number", d.serialNumber).Collect(ch)
	metricWaterMeterPulsarPulsedPulses.With("serial_number", d.serialNumber).Collect(ch)
}

func (d *PulsarPulsedWaterMeter) Ping(_ context.Context) bool {
	_, err := d.provider.Version()
	return err == nil
}

func (d *PulsarPulsedWaterMeter) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.updater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-water-meter-pulsar-pulsed-updater-" + d.serialNumber)

	return []workers.Task{
		taskUpdater,
	}
}

func (d *PulsarPulsedWaterMeter) volume(pulses uint64) float64 {
	return (d.volumeOffset*PulsarPulsedWaterMeterScale + float64(pulses*10)) / PulsarPulsedWaterMeterScale
}

func (d *PulsarPulsedWaterMeter) updater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	pulses, err := d.Pulses(ctx)
	if err != nil {
		return nil, err
	}

	metricWaterMeterPulsarPulsedPulses.With("serial_number", d.serialNumber).Set(float64(pulses))
	metricWaterMeterPulsarPulsedVolume.With("serial_number", d.serialNumber).Set(d.volume(pulses))

	return nil, nil
}
