package internal

import (
	"encoding/hex"

	"github.com/dhickie/go-lgtv/control"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/boggart/components/boggart/protocols/apcupsd"
	"github.com/kihamo/boggart/components/boggart/protocols/apcupsd/file"
	"github.com/kihamo/boggart/components/boggart/protocols/apcupsd/nis"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
	"github.com/kihamo/boggart/components/boggart/providers/mobile"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/boggart/components/boggart/providers/samsung/tv"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
	"github.com/vikstrous/zengge-lightcontrol/local"
	"github.com/vikstrous/zengge-lightcontrol/manage"
)

func (c *Component) initVideoRecorders() {
	isapi := hikvision.NewISAPI(
		c.config.String(boggart.ConfigVideoRecorderHikVisionHost),
		c.config.Int64(boggart.ConfigVideoRecorderHikVisionPort),
		c.config.String(boggart.ConfigVideoRecorderHikVisionUsername),
		c.config.String(boggart.ConfigVideoRecorderHikVisionPassword))

	device := devices.NewVideoRecorderHikVision(isapi, c.config.Duration(boggart.ConfigVideoRecorderHikVisionRepeatInterval))
	if c.config.Bool(boggart.ConfigVideoRecorderHikVisionEnabled) {
		device.Enable()
	} else {
		device.Disable()
	}

	c.devicesManager.RegisterWithID(boggart.DeviceIdVideoRecorder.String(), device)
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

	c.devicesManager.Register(device)
}

func (c *Component) initCameras() {
	isapi := hikvision.NewISAPI(
		c.config.String(boggart.ConfigCameraHikVisionHallHost),
		c.config.Int64(boggart.ConfigCameraHikVisionHallPort),
		c.config.String(boggart.ConfigCameraHikVisionHallUsername),
		c.config.String(boggart.ConfigCameraHikVisionHallPassword))

	device := devices.NewCameraHikVision(
		isapi,
		c.config.Uint64(boggart.ConfigCameraHikVisionHallStreamingChannel),
		c.config.Duration(boggart.ConfigCameraHikVisionHallRepeatInterval))
	device.SetDescription(device.Description() + " on hall")

	if c.config.Bool(boggart.ConfigCameraHikVisionHallEnabled) {
		device.Enable()
	} else {
		device.Disable()
	}

	c.devicesManager.RegisterWithID(boggart.DeviceIdCameraHall.String(), device)

	isapi = hikvision.NewISAPI(
		c.config.String(boggart.ConfigCameraHikVisionStreetHost),
		c.config.Int64(boggart.ConfigCameraHikVisionStreetPort),
		c.config.String(boggart.ConfigCameraHikVisionStreetUsername),
		c.config.String(boggart.ConfigCameraHikVisionStreetPassword))

	device = devices.NewCameraHikVision(
		isapi,
		c.config.Uint64(boggart.ConfigCameraHikVisionStreetStreamingChannel),
		c.config.Duration(boggart.ConfigCameraHikVisionStreetRepeatInterval))
	device.SetDescription(device.Description() + " on street")

	if c.config.Bool(boggart.ConfigCameraHikVisionStreetEnabled) {
		device.Enable()
	} else {
		device.Disable()
	}

	c.devicesManager.RegisterWithID(boggart.DeviceIdCameraStreet.String(), device)
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

	c.devicesManager.RegisterWithID(boggart.DeviceIdPhone.String(), device)
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

	device := devices.NewMikrotikRouter(api, c.config.Duration(boggart.ConfigMikrotikRepeatInterval))

	if c.config.Bool(boggart.ConfigMikrotikEnabled) {
		device.Enable()
	} else {
		device.Disable()
	}

	c.devicesManager.RegisterWithID(boggart.DeviceIdRouter.String(), device)
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

	c.devicesManager.RegisterWithID(boggart.DeviceIdElectricityMeter.String(), device)
}

func (c *Component) initGPIO() {
	device := devices.NewDoorGPIOReedSwitch(c.config.Int64(boggart.ConfigDoorsEntrancePin))
	device.SetDescription("Entrance door")

	if c.config.Bool(boggart.ConfigDoorsEnabled) {
		device.Enable()
	} else {
		device.Disable()
	}

	c.devicesManager.RegisterWithID(boggart.DeviceIdEntranceDoor.String(), device)
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

	c.devicesManager.RegisterWithID(boggart.DeviceIdHeatMeter.String(), deviceHeatMeter)

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
	c.devicesManager.RegisterWithID(boggart.DeviceIdWaterMeterCold.String(), deviceWaterMeterCold)

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
	c.devicesManager.RegisterWithID(boggart.DeviceIdWaterMeterHot.String(), deviceWaterMeterHot)
}

func (c *Component) initUPS() {
	var client *apcupsd.Client

	if address := c.config.String(boggart.ConfigApcupsdNISAddress); address != "" {
		client = apcupsd.NewClient(
			nis.NewStatusReader(address),
			nis.NewEventsReader(address))
	} else {
		client = apcupsd.NewClient(
			file.NewStatusReader(c.config.String(boggart.ConfigApcupsdFileStatus)),
			file.NewEventsReader(c.config.String(boggart.ConfigApcupsdFileEvents)))
	}

	device := devices.NewApcupsdUPS(client, c.config.Duration(boggart.ConfigApcupsdRepeatInterval))

	if c.config.Bool(boggart.ConfigApcupsdEnabled) {
		device.Enable()
	} else {
		device.Disable()
	}

	c.devicesManager.RegisterWithID(boggart.DeviceIdUPS.String(), device)
}

func (c *Component) initTV() {
	// Samsung
	api := tv.NewApiV2("192.168.88.169")

	deviceSamsung := devices.NewSamsungTV(api)
	deviceSamsung.SetMac("9C:8C:6E:CF:3F:EE")
	c.devicesManager.RegisterWithID(boggart.DeviceIdTVBedroom.String(), deviceSamsung)

	// LG
	ctrl, err := control.NewTV("192.168.88.162", "", "")
	if err != nil {
		c.logger.Error("Init control of LG TV failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	ctrl.ClientKey = "3d097186c0fb4abf96b03d251a573613"

	deviceLG := devices.NewLGTV(ctrl)
	deviceLG.SetMac("38:8C:50:68:FB:F4")
	c.devicesManager.RegisterWithID(boggart.DeviceIdTVLivingRoom.String(), deviceLG)
}

func (c *Component) initLight() {
	transport, err := local.NewTransport("192.168.88.165:5577")
	if err != nil {
		c.logger.Error("Init transport of Zengge light failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	manager, err := manage.NewManager("192.168.88.165:48899")
	if err != nil {
		c.logger.Error("Init manager of Zengge light failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	device := devices.NewZenggeLight(transport, manager)
	c.devicesManager.RegisterWithID(boggart.DeviceIdLight.String(), device)
}

func (c *Component) initPC() {
	deviceAi := devices.NewDesktopPC()
	deviceAi.SetMac("50:AF:73:24:B6:7A")

	c.devicesManager.Register(deviceAi)

	deviceRPi := devices.NewDesktopPC()
	c.devicesManager.Register(deviceRPi)
}
