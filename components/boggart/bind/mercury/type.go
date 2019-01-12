package mercury

import (
	"math"
	"time"

	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
	"github.com/kihamo/boggart/components/boggart/providers/mercury"

	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*Config)

	var err error
	timeout := time.Second

	if config.RS485.Timeout != "" {
		timeout, err = time.ParseDuration(config.RS485.Timeout)
		if err != nil {
			return nil, err
		}
	}

	provider := mercury.NewElectricityMeter200(
		mercury.ConvertSerialNumber(config.Address),
		rs485.GetConnection(config.RS485.Address, timeout))

	device := &Bind{
		provider: provider,

		tariff1:        math.MaxUint64,
		tariff2:        math.MaxUint64,
		tariff3:        math.MaxUint64,
		tariff4:        math.MaxUint64,
		voltage:        math.MaxUint64,
		amperage:       math.MaxUint64,
		power:          math.MaxInt64,
		batteryVoltage: math.MaxUint64,
	}
	device.Init()

	// TODO: read real serial number
	device.SetSerialNumber(config.Address)

	// TODO: MQTT publish version

	return device, nil
}
