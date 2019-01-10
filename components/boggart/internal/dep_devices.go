package internal

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

func (c *Component) initElectricityMeters() {
	address := c.config.String(boggart.ConfigMercuryDeviceAddress)
	if address == "" {
		return
	}

	provider := mercury.NewElectricityMeter200(mercury.ConvertSerialNumber(address), c.RS485())

	device := devices.NewMercury200ElectricityMeter(
		c.config.String(boggart.ConfigMercuryDeviceAddress),
		provider,
		c.config.Duration(boggart.ConfigMercuryRepeatInterval))

	c.devicesManager.RegisterWithID(boggart.DeviceIdElectricityMeter.String(), device, "mercury", "", nil, nil)
}

func (c *Component) initGPIO() {
	addresses := c.config.String(boggart.ConfigGPIOPins)
	if addresses == "" {
		return
	}

	for _, pin := range strings.Split(addresses, ",") {
		opts := strings.Split(pin, ":")

		number, err := strconv.ParseUint(opts[0], 10, 64)
		if err != nil {
			continue
		}

		g := gpioreg.ByName(fmt.Sprintf("GPIO%d", number))
		if g == nil {
			c.logger.Warnf("GPIO %d not found", number)
			continue
		}

		var mode devices.GPIOMode
		if len(opts) > 1 {
			switch opts[1] {
			case "in":
				mode = devices.GPIOModeIn
			case "out":
				mode = devices.GPIOModeOut
			default:
				mode = devices.GPIOModeDefault
			}
		}

		device := devices.NewGPIOPin(g, mode)
		c.devicesManager.RegisterWithID(fmt.Sprintf("pin.%d", number), device, "gpio", "", nil, nil)
	}
}

func (c *Component) initPulsarMeters() {
	if !c.config.Bool(boggart.ConfigPulsarEnabled) {
		return
	}

	var (
		deviceAddress []byte
		err           error
	)

	deviceAddressConfig := c.config.String(boggart.ConfigPulsarHeatMeterAddress)
	if deviceAddressConfig == "" {
		deviceAddress, err = pulsar.DeviceAddress(c.RS485())
	} else {
		deviceAddress, err = hex.DecodeString(deviceAddressConfig)
	}

	if err != nil {
		c.logger.Error("Try to get pulsar heat meter address failed", "error", err.Error())
		return
	}

	if len(deviceAddress) != 4 {
		c.logger.Error("Length of device address of pulsar heat meter is wrong")
		return
	}

	provider := pulsar.NewHeatMeter(deviceAddress, c.RS485())

	// heat meter
	deviceHeatMeter := devices.NewPulsarHeadMeter(provider, c.config.Duration(boggart.ConfigPulsarRepeatInterval))

	c.devicesManager.RegisterWithID(boggart.DeviceIdHeatMeter.String(), deviceHeatMeter, "pulsar", "", nil, nil)

	// cold water
	serialNumber := c.config.String(boggart.ConfigPulsarColdWaterSerialNumber)
	deviceWaterMeterCold := devices.NewPulsarPulsedWaterMeter(
		serialNumber,
		c.config.Float64(boggart.ConfigPulsarColdWaterStartValue),
		provider,
		c.config.Uint64(boggart.ConfigPulsarColdWaterPulseInput),
		c.config.Duration(boggart.ConfigPulsarRepeatInterval))

	c.devicesManager.RegisterWithID(boggart.DeviceIdWaterMeterCold.String(), deviceWaterMeterCold, "pulsar", "", nil, nil)

	// hot water
	serialNumber = c.config.String(boggart.ConfigPulsarHotWaterSerialNumber)
	deviceWaterMeterHot := devices.NewPulsarPulsedWaterMeter(
		c.config.String(boggart.ConfigPulsarHotWaterSerialNumber),
		c.config.Float64(boggart.ConfigPulsarHotWaterStartValue),
		provider,
		c.config.Uint64(boggart.ConfigPulsarHotWaterPulseInput),
		c.config.Duration(boggart.ConfigPulsarRepeatInterval))

	c.devicesManager.RegisterWithID(boggart.DeviceIdWaterMeterHot.String(), deviceWaterMeterHot, "pulsar", "", nil, nil)
}
