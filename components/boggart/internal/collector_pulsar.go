package internal

import (
	"encoding/hex"
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
)

func (c *MetricsCollector) CollectPulsar() error {
	connection := pulsar.NewConnection(
		c.component.config.GetString(boggart.ConfigPulsarSerialAddress),
		c.component.config.GetDuration(boggart.ConfigPulsarSerialTimeout))

	var (
		deviceAddress []byte
		err           error
	)

	deviceAddressConfig := c.component.config.GetString(boggart.ConfigPulsarDeviceAddress)
	if deviceAddressConfig == "" {
		deviceAddress, err = connection.DeviceAddress()
	} else {
		deviceAddress, err = hex.DecodeString(deviceAddressConfig)
	}

	if err != nil {
		return err
	}

	if len(deviceAddress) != 4 {
		return errors.New("Length of device address is wrong")
	}

	device := pulsar.NewDevice(deviceAddress, connection)

	temperatureIn, err := device.TemperatureIn()
	if err != nil {
		return err
	}
	metricPulsarTemperatureIn.Set(float64(temperatureIn))

	temperatureOut, err := device.TemperatureOut()
	if err != nil {
		return err
	}
	metricPulsarTemperatureOut.Set(float64(temperatureOut))

	temperatureDelta, err := device.TemperatureDelta()
	if err != nil {
		return err
	}
	metricPulsarTemperatureDelta.Set(float64(temperatureDelta))

	energy, err := device.Energy()
	if err != nil {
		return err
	}
	metricPulsarEnergy.Set(float64(energy))

	consumption, err := device.Consumption()
	if err != nil {
		return err
	}
	metricPulsarConsumption.Set(float64(consumption))

	var coldWaterCapacityFunc func() (float32, error)
	switch c.component.config.GetInt64(boggart.ConfigPulsarColdWaterPulseInput) {
	case pulsar.Input1:
		coldWaterCapacityFunc = device.PulseInput1
	case pulsar.Input2:
		coldWaterCapacityFunc = device.PulseInput2
	default:
		return errors.New("Unknown input of cold water")
	}

	coldWaterCapacity, err := coldWaterCapacityFunc()
	if err != nil {
		return err
	}
	metricPulsarColdWaterCapacity.Set(
		(c.component.config.GetFloat64(boggart.ConfigPulsarColdWaterStartValue)*1000 +
			float64(coldWaterCapacity*10)) / 1000)

	var hotWaterCapacityFunc func() (float32, error)
	switch c.component.config.GetInt64(boggart.ConfigPulsarHotWaterPulseInput) {
	case pulsar.Input1:
		hotWaterCapacityFunc = device.PulseInput1
	case pulsar.Input2:
		hotWaterCapacityFunc = device.PulseInput2
	default:
		return errors.New("Unknown input of hot water")
	}

	hotWaterCapacity, err := hotWaterCapacityFunc()
	if err != nil {
		return err
	}
	metricPulsarHotWaterCapacity.Set(
		(c.component.config.GetFloat64(boggart.ConfigPulsarHotWaterStartValue)*1000 +
			float64(hotWaterCapacity*10)) / 1000)

	return nil
}
