package internal

import (
	"encoding/hex"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
	"github.com/kihamo/boggart/components/boggart/providers/mobile"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
)

func (c *Component) initVideoRecorders() {
	isapi := hikvision.NewISAPI(
		c.config.String(boggart.ConfigVideoRecorderHikVisionHost),
		c.config.Int64(boggart.ConfigVideoRecorderHikVisionPort),
		c.config.String(boggart.ConfigVideoRecorderHikVisionUsername),
		c.config.String(boggart.ConfigVideoRecorderHikVisionPassword))

	device, err := devices.NewVideoRecorderHikVision(isapi, c.config.Duration(boggart.ConfigVideoRecorderHikVisionRepeatInterval))
	if err != nil {
		c.logger.Error("Init video recorder failed", map[string]interface{}{
			"error":    err.Error(),
			"host":     c.config.String(boggart.ConfigVideoRecorderHikVisionHost),
			"port":     c.config.String(boggart.ConfigVideoRecorderHikVisionPort),
			"username": c.config.String(boggart.ConfigVideoRecorderHikVisionUsername),
		})
		return
	}

	if !c.config.Bool(boggart.ConfigVideoRecorderHikVisionEnabled) {
		device.Enable()
	} else {
		device.Disable()
	}

	c.devices.RegisterWithID(boggart.DeviceIdVideoRecorder.String(), device)
}

func (c *Component) initInternetProviders() {
	provider := softvideo.NewClient(
		c.config.String(boggart.ConfigSoftVideoLogin),
		c.config.String(boggart.ConfigSoftVideoPassword))

	device := devices.NewSoftVideoInternet(provider, c.config.Duration(boggart.ConfigSoftVideoRepeatInterval))

	if c.config.Bool(boggart.ConfigSoftVideoEnabled) {
		device.Enable()
	} else {
		device.Disable()
	}

	c.devices.Register(device)
}

func (c *Component) initCameras() {
	isapi := hikvision.NewISAPI(
		c.config.String(boggart.ConfigCameraHikVisionHallHost),
		c.config.Int64(boggart.ConfigCameraHikVisionHallPort),
		c.config.String(boggart.ConfigCameraHikVisionHallUsername),
		c.config.String(boggart.ConfigCameraHikVisionHallPassword))

	device, err := devices.NewCameraHikVision(
		isapi,
		c.config.Uint64(boggart.ConfigCameraHikVisionHallStreamingChannel),
		c.config.Duration(boggart.ConfigCameraHikVisionHallRepeatInterval))
	if err != nil {
		c.logger.Error("Init camera failed", map[string]interface{}{
			"error":    err.Error(),
			"host":     c.config.String(boggart.ConfigCameraHikVisionHallHost),
			"port":     c.config.String(boggart.ConfigCameraHikVisionHallPort),
			"username": c.config.String(boggart.ConfigCameraHikVisionHallUsername),
		})
		return
	} else {
		if c.config.Bool(boggart.ConfigCameraHikVisionHallEnabled) {
			device.Enable()
		} else {
			device.Disable()
		}

		c.devices.RegisterWithID(boggart.DeviceIdCameraHall.String(), device)
	}

	isapi = hikvision.NewISAPI(
		c.config.String(boggart.ConfigCameraHikVisionStreetHost),
		c.config.Int64(boggart.ConfigCameraHikVisionStreetPort),
		c.config.String(boggart.ConfigCameraHikVisionStreetUsername),
		c.config.String(boggart.ConfigCameraHikVisionStreetPassword))

	device, err = devices.NewCameraHikVision(
		isapi,
		c.config.Uint64(boggart.ConfigCameraHikVisionStreetStreamingChannel),
		c.config.Duration(boggart.ConfigCameraHikVisionStreetRepeatInterval))
	if err != nil {
		c.logger.Error("Init camera failed", map[string]interface{}{
			"error":    err.Error(),
			"host":     c.config.String(boggart.ConfigCameraHikVisionStreetHost),
			"port":     c.config.String(boggart.ConfigCameraHikVisionStreetPort),
			"username": c.config.String(boggart.ConfigCameraHikVisionStreetUsername),
		})
		return
	} else {
		if c.config.Bool(boggart.ConfigCameraHikVisionStreetEnabled) {
			device.Enable()
		} else {
			device.Disable()
		}

		c.devices.RegisterWithID(boggart.DeviceIdCameraStreet.String(), device)
	}
}

func (c *Component) initPhones() {
	megafonPhone := c.config.String(boggart.ConfigMobileMegafonPhone)
	megafonPassword := c.config.String(boggart.ConfigMobileMegafonPassword)

	if megafonPhone == "" || megafonPassword == "" {
		return
	}

	providerMegafon := mobile.NewMegafon(megafonPhone, megafonPassword)
	device := devices.NewMegafonPhone(providerMegafon, c.config.Duration(boggart.ConfigMobileRepeatInterval))

	if c.config.Bool(boggart.ConfigMobileEnabled) {
		device.Enable()
	} else {
		device.Disable()
	}

	c.devices.RegisterWithID(boggart.DeviceIdPhone.String(), device)
}

func (c *Component) initRouters() {
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

	device, err := devices.NewMikrotikRouter(api, c.config.Duration(boggart.ConfigMikrotikRepeatInterval))
	if err != nil {
		c.logger.Error("Init router device failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if c.config.Bool(boggart.ConfigMikrotikEnabled) {
		device.Enable()
	} else {
		device.Disable()
	}

	c.devices.RegisterWithID(boggart.DeviceIdRouter.String(), device)
}

func (c *Component) initElectricityMeters() {
	provider := mercury.NewElectricityMeter200(
		mercury.ConvertSerialNumber(c.config.String(boggart.ConfigMercuryDeviceAddress)),
		c.ConnectionRS485())

	device := devices.NewMercury200ElectricityMeter(
		c.config.String(boggart.ConfigMercuryDeviceAddress),
		provider,
		c.config.Duration(boggart.ConfigMercuryRepeatInterval))

	if c.config.Bool(boggart.ConfigMikrotikEnabled) {
		device.Enable()
	} else {
		device.Disable()
	}

	c.devices.RegisterWithID(boggart.DeviceIdElectricityMeter.String(), device)
}

func (c *Component) initPulsarMeters() {
	var (
		deviceAddress []byte
		err           error
	)

	deviceAddressConfig := c.config.String(boggart.ConfigPulsarHeatMeterAddress)
	if deviceAddressConfig == "" {
		deviceAddress, err = pulsar.DeviceAddress(c.ConnectionRS485())
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

	provider := pulsar.NewHeatMeter(deviceAddress, c.ConnectionRS485())

	// heat meter
	deviceHeatMeter := devices.NewPulsarHeadMeter(provider, c.config.Duration(boggart.ConfigPulsarRepeatInterval))

	if c.config.Bool(boggart.ConfigPulsarEnabled) {
		deviceHeatMeter.Enable()
	} else {
		deviceHeatMeter.Disable()
	}

	c.devices.RegisterWithID(boggart.DeviceIdHeatMeter.String(), deviceHeatMeter)

	// cold water
	serialNumber := c.config.String(boggart.ConfigPulsarColdWaterSerialNumber)
	deviceWaterMeterCold := devices.NewPulsarPulsedWaterMeter(
		serialNumber,
		c.config.Float64(boggart.ConfigPulsarColdWaterStartValue),
		provider,
		c.config.Uint64(boggart.ConfigPulsarColdWaterPulseInput),
		c.config.Duration(boggart.ConfigPulsarRepeatInterval))

	if c.config.Bool(boggart.ConfigPulsarEnabled) {
		deviceWaterMeterCold.Enable()
	} else {
		deviceWaterMeterCold.Disable()
	}

	deviceWaterMeterCold.SetDescription("Pulsar pulsed cold water meter with serial number " + serialNumber)
	c.devices.RegisterWithID(boggart.DeviceIdWaterMeterCold.String(), deviceWaterMeterCold)

	// hot water
	serialNumber = c.config.String(boggart.ConfigPulsarHotWaterSerialNumber)
	deviceWaterMeterHot := devices.NewPulsarPulsedWaterMeter(
		c.config.String(boggart.ConfigPulsarHotWaterSerialNumber),
		c.config.Float64(boggart.ConfigPulsarHotWaterStartValue),
		provider,
		c.config.Uint64(boggart.ConfigPulsarHotWaterPulseInput),
		c.config.Duration(boggart.ConfigPulsarRepeatInterval))

	if c.config.Bool(boggart.ConfigPulsarEnabled) {
		deviceWaterMeterHot.Enable()
	} else {
		deviceWaterMeterHot.Disable()
	}

	deviceWaterMeterHot.SetDescription("Pulsar pulsed hot water meter with serial number " + serialNumber)
	c.devices.RegisterWithID(boggart.DeviceIdWaterMeterHot.String(), deviceWaterMeterHot)
}
