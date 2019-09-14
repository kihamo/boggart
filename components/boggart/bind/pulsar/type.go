package pulsar

import (
	"encoding/hex"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/serial"
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

	sn := hex.EncodeToString(deviceAddress)

	config.TopicTemperatureIn = config.TopicTemperatureIn.Format(sn)
	config.TopicTemperatureOut = config.TopicTemperatureOut.Format(sn)
	config.TopicTemperatureDelta = config.TopicTemperatureDelta.Format(sn)
	config.TopicEnergy = config.TopicEnergy.Format(sn)
	config.TopicConsumption = config.TopicConsumption.Format(sn)
	config.TopicCapacity = config.TopicCapacity.Format(sn)
	config.TopicPower = config.TopicPower.Format(sn)
	config.TopicInputPulses1 = config.TopicInputPulses1.Format(sn)
	config.TopicInputPulses2 = config.TopicInputPulses2.Format(sn)
	config.TopicInputPulses3 = config.TopicInputPulses3.Format(sn)
	config.TopicInputPulses4 = config.TopicInputPulses4.Format(sn)
	config.TopicInputVolume1 = config.TopicInputVolume1.Format(sn)
	config.TopicInputVolume2 = config.TopicInputVolume2.Format(sn)
	config.TopicInputVolume3 = config.TopicInputVolume3.Format(sn)
	config.TopicInputVolume4 = config.TopicInputVolume4.Format(sn)

	bind := &Bind{
		config:   config,
		provider: pulsar.NewHeatMeter(deviceAddress, loc, conn),
	}
	bind.SetSerialNumber(sn)

	return bind, nil
}
