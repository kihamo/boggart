package v1

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/serial"
	mercury "github.com/kihamo/boggart/components/boggart/providers/mercury/v1"
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

	bind := &Bind{
		provider:        provider,
		updaterInterval: config.UpdaterInterval,
	}

	// TODO: read real serial number
	bind.SetSerialNumber(config.Address)

	// TODO: MQTT publish version

	return bind, nil
}
