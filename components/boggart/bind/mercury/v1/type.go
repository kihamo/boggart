package v1

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	mercury "github.com/kihamo/boggart/providers/mercury/v1"
)

type Type struct {
	SerialNumberFunc func(address string) mercury.Option
	Device           uint8
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	config.TopicTariff = config.TopicTariff.Format(config.Address)
	config.TopicVoltage = config.TopicVoltage.Format(config.Address)
	config.TopicAmperage = config.TopicAmperage.Format(config.Address)
	config.TopicPower = config.TopicPower.Format(config.Address)
	config.TopicBatteryVoltage = config.TopicBatteryVoltage.Format(config.Address)
	config.TopicLastPowerOff = config.TopicLastPowerOff.Format(config.Address)
	config.TopicLastPowerOn = config.TopicLastPowerOn.Format(config.Address)
	config.TopicMakeDate = config.TopicMakeDate.Format(config.Address)
	config.TopicFirmwareDate = config.TopicFirmwareDate.Format(config.Address)
	config.TopicFirmwareVersion = config.TopicFirmwareVersion.Format(config.Address)

	return &Bind{
		config:       config,
		providerOnce: &atomic.Once{},
		tariffCount:  atomic.NewUint32Null(),
	}, nil
}
