package internal

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/snitch"
)

const (
	MetricPulsarColdWaterCapacity = boggart.ComponentName + "_pulsar_cold_water_capacity_cubic_metres"
	MetricPulsarHotWaterCapacity  = boggart.ComponentName + "_pulsar_hot_water_capacity_cubic_metres"
)

var (
	metricPulsarColdWaterCapacity = snitch.NewGauge(MetricPulsarColdWaterCapacity, "Pulsar capacity of cold water")
	metricPulsarHotWaterCapacity  = snitch.NewGauge(MetricPulsarHotWaterCapacity, "Pulsar capacity of hot water")
)

func (c *MetricsCollector) UpdaterPulsar() error {
	var (
		deviceAddress []byte
		err           error
	)

	deviceAddressConfig := c.component.config.String(boggart.ConfigPulsarHeatMeterAddress)
	if deviceAddressConfig == "" {
		deviceAddress, err = pulsar.DeviceAddress(c.component.ConnectionRS485())
	} else {
		deviceAddress, err = hex.DecodeString(deviceAddressConfig)
	}

	if err != nil {
		return fmt.Errorf("DeviceAddress error: %s", err.Error())
	}

	if len(deviceAddress) != 4 {
		return errors.New("Length of device address is wrong")
	}

	device := pulsar.NewHeatMeter(deviceAddress, c.component.ConnectionRS485())

	var coldWaterCapacityFunc func() (float32, error)
	switch c.component.config.Uint64(boggart.ConfigPulsarColdWaterPulseInput) {
	case pulsar.Input1:
		coldWaterCapacityFunc = device.PulseInput1
	case pulsar.Input2:
		coldWaterCapacityFunc = device.PulseInput2
	default:
		return errors.New("Unknown input of cold water")
	}

	coldWaterCapacity, err := coldWaterCapacityFunc()
	if err != nil {
		return fmt.Errorf("ColdWaterCapacityFunc error: %s", err.Error())
	}
	metricPulsarColdWaterCapacity.Set(
		(c.component.config.Float64(boggart.ConfigPulsarColdWaterStartValue)*1000 +
			float64(coldWaterCapacity*10)) / 1000)

	var hotWaterCapacityFunc func() (float32, error)
	switch c.component.config.Uint64(boggart.ConfigPulsarHotWaterPulseInput) {
	case pulsar.Input1:
		hotWaterCapacityFunc = device.PulseInput1
	case pulsar.Input2:
		hotWaterCapacityFunc = device.PulseInput2
	default:
		return errors.New("Unknown input of hot water")
	}

	hotWaterCapacity, err := hotWaterCapacityFunc()
	if err != nil {
		return fmt.Errorf("HotWaterCapacityFunc error: %s", err.Error())
	}
	metricPulsarHotWaterCapacity.Set(
		(c.component.config.Float64(boggart.ConfigPulsarHotWaterStartValue)*1000 +
			float64(hotWaterCapacity*10)) / 1000)

	return nil
}

func (c *MetricsCollector) DescribePulsar(ch chan<- *snitch.Description) {
	metricPulsarColdWaterCapacity.Describe(ch)
	metricPulsarHotWaterCapacity.Describe(ch)
}

func (c *MetricsCollector) CollectPulsar(ch chan<- snitch.Metric) {
	metricPulsarColdWaterCapacity.Collect(ch)
	metricPulsarHotWaterCapacity.Collect(ch)
}
