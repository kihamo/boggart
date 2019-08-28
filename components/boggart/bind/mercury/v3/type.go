package v3

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/serial"
	mercury "github.com/kihamo/boggart/components/boggart/providers/mercury/v3"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	provider := mercury.New(serial.Dial(config.RS485Address, serial.WithTimeout(config.RS485Timeout)))
	if config.Address != "" {
		provider = provider.WithAddress(mercury.ConvertSerialNumber(config.Address))
	}

	bind := &Bind{
		provider: provider,
		config:   config,
	}

	return bind, nil
}
