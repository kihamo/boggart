package mercury

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	provider := mercury.NewElectricityMeter200(
		mercury.ConvertSerialNumber(config.Address),
		rs485.GetConnection(config.RS485Address, config.RS485Timeout))

	bind := &Bind{
		provider: provider,

		tariff1:          atomic.NewUint32Null(),
		tariff2:          atomic.NewUint32Null(),
		tariff3:          atomic.NewUint32Null(),
		tariff4:          atomic.NewUint32Null(),
		voltage:          atomic.NewUint32Null(),
		amperage:         atomic.NewFloat32Null(),
		power:            atomic.NewUint32Null(),
		batteryVoltage:   atomic.NewFloat32Null(),
		makeDate:         atomic.NewUint32Null(),
		lastPowerOffDate: atomic.NewUint32Null(),
		lastPowerOnDate:  atomic.NewUint32Null(),
		firmwareDate:     atomic.NewUint32Null(),
		firmwareVersion:  atomic.NewString(),

		updaterInterval: config.UpdaterInterval,
	}

	// TODO: read real serial number
	bind.SetSerialNumber(config.Address)

	// TODO: MQTT publish version

	return bind, nil
}
