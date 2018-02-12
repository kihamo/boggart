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
	MetricPulsarTemperatureIn     = boggart.ComponentName + "_pulsar_temperature_in_celsius"
	MetricPulsarTemperatureOut    = boggart.ComponentName + "_pulsar_temperature_out_celsius"
	MetricPulsarTemperatureDelta  = boggart.ComponentName + "_pulsar_temperature_delta_celsius"
	MetricPulsarEnergy            = boggart.ComponentName + "_pulsar_energy_gigacolories"
	MetricPulsarConsumption       = boggart.ComponentName + "_pulsar_consumption_cubic_metres_per_hour"
	MetricPulsarColdWaterCapacity = boggart.ComponentName + "_pulsar_cold_water_capacity_cubic_metres"
	MetricPulsarHotWaterCapacity  = boggart.ComponentName + "_pulsar_hot_water_capacity_cubic_metres"
)

var (
	metricPulsarTemperatureIn     = snitch.NewGauge(MetricPulsarTemperatureIn, "Pulsar temperature in")
	metricPulsarTemperatureOut    = snitch.NewGauge(MetricPulsarTemperatureOut, "Pulsar temperature out")
	metricPulsarTemperatureDelta  = snitch.NewGauge(MetricPulsarTemperatureDelta, "Pulsar temperature delta")
	metricPulsarEnergy            = snitch.NewGauge(MetricPulsarEnergy, "Pulsar energy")
	metricPulsarConsumption       = snitch.NewGauge(MetricPulsarConsumption, "Pulsar consumption")
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

	temperatureIn, err := device.TemperatureIn()
	if err != nil {
		return fmt.Errorf("TemperatureIn error: %s", err.Error())
	}
	metricPulsarTemperatureIn.Set(float64(temperatureIn))

	temperatureOut, err := device.TemperatureOut()
	if err != nil {
		return fmt.Errorf("TemperatureOut error: %s", err.Error())
	}
	metricPulsarTemperatureOut.Set(float64(temperatureOut))

	temperatureDelta, err := device.TemperatureDelta()
	if err != nil {
		return fmt.Errorf("TemperatureDelta error: %s", err.Error())
	}
	metricPulsarTemperatureDelta.Set(float64(temperatureDelta))

	energy, err := device.Energy()
	if err != nil {
		return fmt.Errorf("Energy error: %s", err.Error())
	}
	metricPulsarEnergy.Set(float64(energy))

	consumption, err := device.Consumption()
	if err != nil {
		return fmt.Errorf("Consumption error: %s", err.Error())
	}
	metricPulsarConsumption.Set(float64(consumption))

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
	metricPulsarTemperatureIn.Describe(ch)
	metricPulsarTemperatureOut.Describe(ch)
	metricPulsarTemperatureDelta.Describe(ch)
	metricPulsarEnergy.Describe(ch)
	metricPulsarConsumption.Describe(ch)
	metricPulsarColdWaterCapacity.Describe(ch)
	metricPulsarHotWaterCapacity.Describe(ch)
}

func (c *MetricsCollector) CollectPulsar(ch chan<- snitch.Metric) {
	metricPulsarTemperatureIn.Collect(ch)
	metricPulsarTemperatureOut.Collect(ch)
	metricPulsarTemperatureDelta.Collect(ch)
	metricPulsarEnergy.Collect(ch)
	metricPulsarConsumption.Collect(ch)
	metricPulsarColdWaterCapacity.Collect(ch)
	metricPulsarHotWaterCapacity.Collect(ch)
}
