package pulsar

import (
	"encoding/hex"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/connection"
	"github.com/kihamo/boggart/providers/pulsar"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	conn, err := connection.New(config.ConnectionDSN)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(config.Location)
	if err != nil {
		return nil, err
	}

	var deviceAddress []byte
	if config.Address == "" {
		deviceAddress, err = pulsar.DeviceAddress(conn)
	} else {
		deviceAddress, err = hex.DecodeString(config.Address)
	}

	if err != nil {
		return nil, err
	}

	address := hex.EncodeToString(deviceAddress)

	config.TopicTemperatureIn = config.TopicTemperatureIn.Format(address)
	config.TopicTemperatureOut = config.TopicTemperatureOut.Format(address)
	config.TopicTemperatureDelta = config.TopicTemperatureDelta.Format(address)
	config.TopicEnergy = config.TopicEnergy.Format(address)
	config.TopicConsumption = config.TopicConsumption.Format(address)
	config.TopicCapacity = config.TopicCapacity.Format(address)
	config.TopicPower = config.TopicPower.Format(address)
	config.TopicInputPulses1 = config.TopicInputPulses1.Format(address)
	config.TopicInputPulses2 = config.TopicInputPulses2.Format(address)
	config.TopicInputPulses3 = config.TopicInputPulses3.Format(address)
	config.TopicInputPulses4 = config.TopicInputPulses4.Format(address)
	config.TopicInputVolume1 = config.TopicInputVolume1.Format(address)
	config.TopicInputVolume2 = config.TopicInputVolume2.Format(address)
	config.TopicInputVolume3 = config.TopicInputVolume3.Format(address)
	config.TopicInputVolume4 = config.TopicInputVolume4.Format(address)

	opts := []pulsar.Option{
		pulsar.WithAddress(deviceAddress),
		pulsar.WithLocation(loc),
	}

	bind := &Bind{
		config:   config,
		provider: pulsar.New(conn, opts...),
		address:  address,
	}

	return bind, nil
}
