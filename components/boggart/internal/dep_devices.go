package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
)

func (c *Component) initElectricityMeters() {
	address := c.config.String(boggart.ConfigMercuryDeviceAddress)
	if address == "" {
		return
	}

	provider := mercury.NewElectricityMeter200(mercury.ConvertSerialNumber(address), c.RS485())

	device := bind.NewMercury200ElectricityMeter(
		c.config.String(boggart.ConfigMercuryDeviceAddress),
		provider,
		c.config.Duration(boggart.ConfigMercuryRepeatInterval))

	c.devicesManager.RegisterWithID(boggart.DeviceIdElectricityMeter.String(), device, "mercury", "", nil, nil)
}
