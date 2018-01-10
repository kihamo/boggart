package internal

import (
	"encoding/hex"
	"errors"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
	"github.com/kihamo/snitch"
)

type MetricsCollector struct {
	component *Component
}

func NewMetricsCollector(component *Component) *MetricsCollector {
	return &MetricsCollector{
		component: component,
	}
}

func (c *MetricsCollector) Describe(ch chan<- *snitch.Description) {
	metricPulsarTemperatureIn.Describe(ch)
	metricPulsarTemperatureOut.Describe(ch)
	metricPulsarTemperatureDelta.Describe(ch)
	metricSoftVideoBalance.Describe(ch)
}

func (c *MetricsCollector) Collect(ch chan<- snitch.Metric) {
	metricPulsarTemperatureIn.Collect(ch)
	metricPulsarTemperatureOut.Collect(ch)
	metricPulsarTemperatureDelta.Collect(ch)
	metricSoftVideoBalance.Collect(ch)
}

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

	return nil
}

func (c *MetricsCollector) CollectSoftVideo() error {
	client := softvideo.NewClient(
		c.component.config.GetString(boggart.ConfigSoftVideoLogin),
		c.component.config.GetString(boggart.ConfigSoftVideoPassword))

	value, err := client.Balance()
	if err != nil {
		// FIXME: logging
		return err
	}

	metricSoftVideoBalance.Set(float64(value))

	return nil
}
