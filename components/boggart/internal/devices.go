package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
	"github.com/kihamo/boggart/components/boggart/providers/mikrotik"
	"github.com/kihamo/boggart/components/boggart/providers/mobile"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
)

func (c *Component) initVideoRecorders() {
	isapi := hikvision.NewISAPI(
		c.config.String(boggart.ConfigVideoRecorderHikVisionHost),
		c.config.Int64(boggart.ConfigVideoRecorderHikVisionPort),
		c.config.String(boggart.ConfigVideoRecorderHikVisionUsername),
		c.config.String(boggart.ConfigVideoRecorderHikVisionPassword))

	device, err := devices.NewVideoRecorderHikVision(isapi)
	if err != nil {
		c.logger.Error("Init video recorder failed", map[string]interface{}{
			"error":    err.Error(),
			"host":     c.config.String(boggart.ConfigVideoRecorderHikVisionHost),
			"port":     c.config.String(boggart.ConfigVideoRecorderHikVisionPort),
			"username": c.config.String(boggart.ConfigVideoRecorderHikVisionUsername),
		})
		return
	}

	if c.config.Bool(boggart.ConfigVideoRecorderHikVisionEnabled) {
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

	device, err := devices.NewCameraHikVision(isapi, c.config.Uint64(boggart.ConfigCameraHikVisionHallStreamingChannel))
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

	device, err = devices.NewCameraHikVision(isapi, c.config.Uint64(boggart.ConfigCameraHikVisionStreetStreamingChannel))
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
