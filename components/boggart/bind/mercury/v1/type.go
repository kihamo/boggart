package v1

import (
	"net/url"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/connection"
	"github.com/kihamo/boggart/protocols/serial"
	"github.com/kihamo/boggart/protocols/serial_network"
	mercury "github.com/kihamo/boggart/providers/mercury/v1"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	u, err := url.Parse(config.RS485Address)
	if err != nil {
		return nil, err
	}

	var conn connection.Conn

	switch u.Scheme {
	case "tcp", "tcp4", "tcp6":
		conn = serial_network.NewTCPClient(u.Scheme, u.Host)

	case "udp", "udp4", "udp6", "unixgram":
		conn = serial_network.NewUDPClient(u.Scheme, u.Host)

	default:
		conn = connection.NewIO(
			serial.Dial(serial.WithTarget(config.RS485Address),
				serial.WithTimeout(config.RS485Timeout),
			))
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
		mercury.WithAddressAsString(config.Address),
		mercury.WithLocation(loc),
	}

	bind := &Bind{
		config:   config,
		provider: mercury.New(conn, opts...),
	}

	// TODO: read real serial number
	bind.SetSerialNumber(config.Address)

	// TODO: MQTT publish version

	return bind, nil
}
