package devices

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

const (
	PulsarPulsedWaterMeterScale = 1000

	PulsarPulsedWaterMeterMQTTTopicPulses boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/pulsar/+/pulses"
	PulsarPulsedWaterMeterMQTTTopicVolume boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/pulsar/+/volume"
)

var (
	metricWaterMeterPulsarPulsedVolume = snitch.NewGauge(boggart.ComponentName+"_device_water_meter_pulsar_pulsed_volume_cubic_metres", "Pulsar volume of water in cubic metres")
	metricWaterMeterPulsarPulsedPulses = snitch.NewGauge(boggart.ComponentName+"_device_water_meter_pulsar_pulsed_pulses", "Pulsar volume of water in pulses")
)

type PulsarPulsedWaterMeter struct {
	pulses uint64

	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	input        uint64
	volumeOffset float64
	provider     *pulsar.HeatMeter
	interval     time.Duration
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
	volume := d.volume(pulses)

	metricWaterMeterPulsarPulsedPulses.With("serial_number", serialNumber).Set(float64(pulses))
	metricWaterMeterPulsarPulsedVolume.With("serial_number", serialNumber).Set(volume)

	prevPulses := atomic.LoadUint64(&d.pulses)
	if pulses != prevPulses {
		atomic.StoreUint64(&d.pulses, pulses)

		d.MQTTPublishAsync(ctx, PulsarPulsedWaterMeterMQTTTopicPulses.Format(serialNumber), 0, true, pulses)
		d.MQTTPublishAsync(ctx, PulsarPulsedWaterMeterMQTTTopicVolume.Format(serialNumber), 0, true, volume)
	}

	return nil, nil
}

func (d *PulsarPulsedWaterMeter) MQTTTopics() []boggart.DeviceMQTTTopic {
	return []boggart.DeviceMQTTTopic{
		PulsarPulsedWaterMeterMQTTTopicPulses,
		PulsarPulsedWaterMeterMQTTTopicVolume,
	}
}
