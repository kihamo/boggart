package pulsar

import (
	"encoding/hex"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/serial"
	"github.com/kihamo/boggart/providers/pulsar"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	var err error
	conn := serial.Dial(config.RS485Address, serial.WithTimeout(config.RS485Timeout))

	var deviceAddress []byte
	if config.Address == "" {
		deviceAddress, err = pulsar.DeviceAddress(conn)
	} else {
		deviceAddress, err = hex.DecodeString(config.Address)
	}

	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(config.Location)
	if err != nil {
		return nil, err
	}

	bind := &Bind{
		config:          config,
		provider:        pulsar.NewHeatMeter(deviceAddress, loc, conn),
		updaterInterval: config.UpdaterInterval,
	}
	bind.SetSerialNumber(hex.EncodeToString(deviceAddress))

	return bind, nil
}
