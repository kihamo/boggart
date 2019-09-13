package v1

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/serial"
	mercury "github.com/kihamo/boggart/providers/mercury/v1"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	loc, err := time.LoadLocation(config.Location)
	if err != nil {
		return nil, err
	}

	provider := mercury.New(
		mercury.ConvertSerialNumber(config.Address),
		loc,
		serial.Dial(config.RS485Address, serial.WithTimeout(config.RS485Timeout)))

	config.TopicTariff1 = config.TopicTariff1.Format(config.Address)
	config.TopicTariff2 = config.TopicTariff2.Format(config.Address)
	config.TopicTariff3 = config.TopicTariff3.Format(config.Address)
	config.TopicTariff4 = config.TopicTariff4.Format(config.Address)
	config.TopicVoltage = config.TopicVoltage.Format(config.Address)
	config.TopicAmperage = config.TopicAmperage.Format(config.Address)
	config.TopicPower = config.TopicPower.Format(config.Address)
	config.TopicBatteryVoltage = config.TopicBatteryVoltage.Format(config.Address)
	config.TopicLastPowerOff = config.TopicLastPowerOff.Format(config.Address)
	config.TopicLastPowerOn = config.TopicLastPowerOn.Format(config.Address)
	config.TopicMakeDate = config.TopicMakeDate.Format(config.Address)
	config.TopicFirmwareDate = config.TopicFirmwareDate.Format(config.Address)
	config.TopicFirmwareVersion = config.TopicFirmwareVersion.Format(config.Address)

	bind := &Bind{
		config:   config,
		provider: provider,
	}

	// TODO: read real serial number
	bind.SetSerialNumber(config.Address)

	// TODO: MQTT publish version

	return bind, nil
}
