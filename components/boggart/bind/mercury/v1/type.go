package v1

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/connection"
	mercury "github.com/kihamo/boggart/providers/mercury/v1"
)

type Type struct {
	boggart.BindTypeWidget

	SerialNumberFunc func(address string) mercury.Option
	Device           uint8
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	conn, err := connection.NewByDSNString(config.ConnectionDSN)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(config.Location)
	if err != nil {
		return nil, err
	}

	config.TopicTariff1 = config.TopicTariff1.Format(config.Address)
	config.TopicTariff2 = config.TopicTariff2.Format(config.Address)
	config.TopicTariff3 = config.TopicTariff3.Format(config.Address)
	config.TopicTariff4 = config.TopicTariff4.Format(config.Address)
	config.TopicVoltage = config.TopicVoltage.Format(config.Address)
	config.TopicAmperage = config.TopicAmperage.Format(config.Address)
	config.TopicPower = config.TopicPower.Format(config.Address)
	config.TopicBatteryVoltage = config.TopicBatteryVoltage.Format(config.Address)
	config.TopicLastPowerOff = config.TopicLastPowerOff.Format(config.Address)
	config.TopicLastPowerOn = config.TopicLastPowerOn.Format(config.Address)
	config.TopicMakeDate = config.TopicMakeDate.Format(config.Address)
	config.TopicFirmwareDate = config.TopicFirmwareDate.Format(config.Address)
	config.TopicFirmwareVersion = config.TopicFirmwareVersion.Format(config.Address)

	opts := []mercury.Option{
		t.SerialNumberFunc(config.Address),
		mercury.WithDevice(t.Device),
		mercury.WithLocation(loc),
	}

	bind := &Bind{
		config:   config,
		provider: mercury.New(conn, opts...),
	}

	conn.ApplyOptions(connection.WithDumpRead(func(bytes []byte) {
		bind.Logger().Debug("Read packet",
			"payload", fmt.Sprintf("%v", bytes),
			"hex", hex.EncodeToString(bytes),
		)
	}))
	conn.ApplyOptions(connection.WithDumpWrite(func(bytes []byte) {
		bind.Logger().Debug("Write packet",
			"payload", fmt.Sprintf("%v", bytes),
			"hex", hex.EncodeToString(bytes),
		)
	}))

	// TODO: MQTT publish version

	return bind, nil
}
