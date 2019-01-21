package pulsar

import (
	"encoding/hex"
	"math"

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

		temperatureIn:    atomic.NewFloat32Default(math.MaxFloat32),
		temperatureOut:   atomic.NewFloat32Default(math.MaxFloat32),
		temperatureDelta: atomic.NewFloat32Default(math.MaxFloat32),
		energy:           atomic.NewFloat32Default(math.MaxFloat32),
		consumption:      atomic.NewFloat32Default(math.MaxFloat32),
		capacity:         atomic.NewFloat32Default(math.MaxFloat32),
		power:            atomic.NewFloat32Default(math.MaxFloat32),
		input1:           atomic.NewFloat32Default(math.MaxFloat32),
		input2:           atomic.NewFloat32Default(math.MaxFloat32),
		input3:           atomic.NewFloat32Default(math.MaxFloat32),
		input4:           atomic.NewFloat32Default(math.MaxFloat32),

		updaterInterval: config.UpdaterInterval,
	}
	device.SetSerialNumber(hex.EncodeToString(deviceAddress))

	return device, nil
}
