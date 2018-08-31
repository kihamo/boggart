package internal

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
	"github.com/kihamo/boggart/components/boggart/providers/mobile"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
	"github.com/yryz/ds18b20"
	"gobot.io/x/gobot/platforms/raspi"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

func (c *Component) initVideoRecorders() {
	var (
		isapi  *hikvision.ISAPI
		device *devices.VideoRecorderHikVision
	)

	if c.config.Bool(boggart.ConfigVideoRecorderHikVisionHomeEnabled) {
		isapi = hikvision.NewISAPI(
			c.config.String(boggart.ConfigVideoRecorderHikVisionHomeHost),
			c.config.Int64(boggart.ConfigVideoRecorderHikVisionHomePort),
			c.config.String(boggart.ConfigVideoRecorderHikVisionHomeUsername),
			c.config.String(boggart.ConfigVideoRecorderHikVisionHomePassword))

		device = devices.NewVideoRecorderHikVision(isapi, c.config.Duration(boggart.ConfigVideoRecorderHikVisionHomeRepeatInterval))
		device.SetDescription("Home video recorder")

		c.devicesManager.Register(device)
	}

	if c.config.Bool(boggart.ConfigVideoRecorderHikVisionVacationHomeEnabled) {
		isapi = hikvision.NewISAPI(
			c.config.String(boggart.ConfigVideoRecorderHikVisionVacationHomeHost),
			c.config.Int64(boggart.ConfigVideoRecorderHikVisionVacationHomePort),
			c.config.String(boggart.ConfigVideoRecorderHikVisionVacationHomeUsername),
			c.config.String(boggart.ConfigVideoRecorderHikVisionVacationHomePassword))

		device = devices.NewVideoRecorderHikVision(isapi, c.config.Duration(boggart.ConfigVideoRecorderHikVisionVacationHomeRepeatInterval))
		device.SetDescription("Vacation home video recorder")

		c.devicesManager.Register(device)
	}

	if c.config.Bool(boggart.ConfigVideoRecorderHikVisionGarageEnabled) {
		isapi = hikvision.NewISAPI(
			c.config.String(boggart.ConfigVideoRecorderHikVisionGarageHost),
			c.config.Int64(boggart.ConfigVideoRecorderHikVisionGaragePort),
			c.config.String(boggart.ConfigVideoRecorderHikVisionGarageUsername),
			c.config.String(boggart.ConfigVideoRecorderHikVisionGaragePassword))

		device = devices.NewVideoRecorderHikVision(isapi, c.config.Duration(boggart.ConfigVideoRecorderHikVisionGarageRepeatInterval))
		device.SetDescription("Garage video recorder")

		c.devicesManager.Register(device)
	}
}

func (c *Component) initInternetProviders() {
	if !c.config.Bool(boggart.ConfigSoftVideoEnabled) {
		return
	}

	provider := softvideo.NewClient(
		c.config.String(boggart.ConfigSoftVideoLogin),
		c.config.String(boggart.ConfigSoftVideoPassword))

	device := devices.NewSoftVideoInternet(provider, c.config.Duration(boggart.ConfigSoftVideoRepeatInterval))

	c.devicesManager.Register(device)
}

func (c *Component) initPhones() {
	if !c.config.Bool(boggart.ConfigMobileEnabled) {
		return
	}

	megafonPhone := c.config.String(boggart.ConfigMobileMegafonPhone)
	megafonPassword := c.config.String(boggart.ConfigMobileMegafonPassword)

	if megafonPhone == "" || megafonPassword == "" {
		return
	}

	providerMegafon := mobile.NewMegafon(megafonPhone, megafonPassword)
	device := devices.NewMegafonPhone(providerMegafon, c.config.Duration(boggart.ConfigMobileRepeatInterval))

	c.devicesManager.RegisterWithID(boggart.DeviceIdPhone.String(), device)
}

func (c *Component) initRouters() {
	if !c.config.Bool(boggart.ConfigMikrotikEnabled) {
		return
	}

	api, err := mikrotik.NewClient(
		c.config.String(boggart.ConfigMikrotikAddress),
		c.config.String(boggart.ConfigMikrotikUsername),
		c.config.String(boggart.ConfigMikrotikPassword),
		c.config.Duration(boggart.ConfigMikrotikTimeout))
	if err != nil {
		c.logger.Error("Init mikrotik api failed", map[string]interface{}{
			"error":    err.Error(),
			"address":  c.config.String(boggart.ConfigMikrotikAddress),
			"username": c.config.String(boggart.ConfigMikrotikUsername),
		})
		return
	}

	device := devices.NewMikrotikRouter(api, c.config.Duration(boggart.ConfigMikrotikRepeatInterval))

	c.devicesManager.RegisterWithID(boggart.DeviceIdRouter.String(), device)
}

func (c *Component) initElectricityMeters() {
	if !c.config.Bool(boggart.ConfigMercuryEnabled) {
		return
	}

	provider := mercury.NewElectricityMeter200(
		mercury.ConvertSerialNumber(c.config.String(boggart.ConfigMercuryDeviceAddress)),
		c.RS485())

	device := devices.NewMercury200ElectricityMeter(
		c.config.String(boggart.ConfigMercuryDeviceAddress),
		provider,
		c.config.Duration(boggart.ConfigMercuryRepeatInterval))

	c.devicesManager.RegisterWithID(boggart.DeviceIdElectricityMeter.String(), device)
}

func (c *Component) initGPIO() {
	if !c.config.Bool(boggart.ConfigGPIOEnabled) {
		return
	}

	pins := strings.Split(c.config.String(boggart.ConfigGPIOPins), ",")

	for _, pin := range pins {
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

		if len(opts) > 2 {
			device.SetDescription(opts[2])
		}

		if c.config.Bool(boggart.ConfigGPIOEnabled) {
			device.Enable()
		} else {
			device.Disable()
		}

		c.devicesManager.RegisterWithID(fmt.Sprintf("pin.%d", number), device)
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
		c.logger.Error("Try to get pulsar heat meter address failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if len(deviceAddress) != 4 {
		c.logger.Error("Length of device address of pulsar heat meter is wrong")
		return
	}

	provider := pulsar.NewHeatMeter(deviceAddress, c.RS485())

	// heat meter
	deviceHeatMeter := devices.NewPulsarHeadMeter(provider, c.config.Duration(boggart.ConfigPulsarRepeatInterval))

	c.devicesManager.RegisterWithID(boggart.DeviceIdHeatMeter.String(), deviceHeatMeter)

	// cold water
	serialNumber := c.config.String(boggart.ConfigPulsarColdWaterSerialNumber)
	deviceWaterMeterCold := devices.NewPulsarPulsedWaterMeter(
		serialNumber,
		c.config.Float64(boggart.ConfigPulsarColdWaterStartValue),
		provider,
		c.config.Uint64(boggart.ConfigPulsarColdWaterPulseInput),
		c.config.Duration(boggart.ConfigPulsarRepeatInterval))

	deviceWaterMeterCold.SetDescription("Pulsar pulsed cold water meter with serial number " + serialNumber)
	c.devicesManager.RegisterWithID(boggart.DeviceIdWaterMeterCold.String(), deviceWaterMeterCold)

	// hot water
	serialNumber = c.config.String(boggart.ConfigPulsarHotWaterSerialNumber)
	deviceWaterMeterHot := devices.NewPulsarPulsedWaterMeter(
		c.config.String(boggart.ConfigPulsarHotWaterSerialNumber),
		c.config.Float64(boggart.ConfigPulsarHotWaterStartValue),
		provider,
		c.config.Uint64(boggart.ConfigPulsarHotWaterPulseInput),
		c.config.Duration(boggart.ConfigPulsarRepeatInterval))

	deviceWaterMeterHot.SetDescription("Pulsar pulsed hot water meter with serial number " + serialNumber)
	c.devicesManager.RegisterWithID(boggart.DeviceIdWaterMeterHot.String(), deviceWaterMeterHot)
}

func (c *Component) initSensor() {
	if c.config.Bool(boggart.ConfigSensorBME280Enabled) {
		deviceBME280 := devices.NewBME280Sensor(
			raspi.NewAdaptor(),
			c.config.Duration(boggart.ConfigSensorBME280RepeatInterval),
			c.config.Int(boggart.ConfigSensorBME280Bus),
			c.config.Int(boggart.ConfigSensorBME280Address))

		c.devicesManager.Register(deviceBME280)
	}

	sensors, err := ds18b20.Sensors()
	if err == nil && len(sensors) > 0 && sensors[0] != "not found." {
		for _, sensor := range sensors {
			device := devices.NewDS18B20Sensor(sensor)
			c.devicesManager.Register(device)
		}
	}
}
