package neptun

import (
	"net/url"

	"github.com/kihamo/boggart/protocols/modbus"
)

const (
	AddressModuleConfiguration       uint16 = 0
	AddressInputLines12Configuration uint16 = 1
	AddressInputLines34Configuration uint16 = 2
	AddressInputLinesStatus          uint16 = 3
	AddressEventsRelayConfiguration  uint16 = 4
	AddressSlaveIDAndBaudRate        uint16 = 5
	AddressWirelessSensorCount       uint16 = 6
	Counter1Slot1Configuration       uint16 = 123
	Counter2Slot1Configuration       uint16 = 124
	Counter1Slot2Configuration       uint16 = 125
	Counter2Slot2Configuration       uint16 = 126
	Counter1Slot3Configuration       uint16 = 127
	Counter2Slot3Configuration       uint16 = 128
	Counter1Slot4Configuration       uint16 = 129
	Counter2Slot4Configuration       uint16 = 130
)

type Neptun struct {
	client *modbus.Client
}

func New(address *url.URL, opts ...modbus.Option) *Neptun {
	opts = append([]modbus.Option{
		modbus.WithSlaveID(240),
	}, opts...)

	address.Scheme = "tcp"

	return &Neptun{
		client: modbus.NewClient(address, opts...),
	}
}
