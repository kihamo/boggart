package mercury

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/protocols/serial"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"
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

	provider := mercury.NewMercury(
		mercury.ConvertSerialNumber(config.Address),
		loc,
		serial.Dial(config.RS485Address, serial.WithTimeout(config.RS485Timeout)))

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
