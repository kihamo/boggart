package pulsar

import (
	"encoding/hex"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	var err error
	conn := rs485.GetConnection(config.RS485Address, config.RS485Timeout)

	var deviceAddress []byte
	if config.Address == "" {
		deviceAddress, err = pulsar.DeviceAddress(conn)
	} else {
		deviceAddress, err = hex.DecodeString(config.Address)
	}

	if err != nil {
		return nil, err
	}

	device := &Bind{
		config:   config,
		provider: pulsar.NewHeatMeter(deviceAddress, conn),

		temperatureIn:    atomic.NewFloat32Null(),
		temperatureOut:   atomic.NewFloat32Null(),
		temperatureDelta: atomic.NewFloat32Null(),
		energy:           atomic.NewFloat32Null(),
		consumption:      atomic.NewFloat32Null(),
		capacity:         atomic.NewFloat32Null(),
		power:            atomic.NewFloat32Null(),
		input1:           atomic.NewFloat32Null(),
		input2:           atomic.NewFloat32Null(),
		input3:           atomic.NewFloat32Null(),
		input4:           atomic.NewFloat32Null(),

		updaterInterval: config.UpdaterInterval,
	}
	device.SetSerialNumber(hex.EncodeToString(deviceAddress))

	return device, nil
}
