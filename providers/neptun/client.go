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
	AddressCounter1ValueHigh         uint16 = 107
	AddressCounter1ValueLow          uint16 = 108
)

const (
	EventsRelaySwitchNo = uint8(iota)
	EventsRelaySwitchFirstGroup
	EventsRelaySwitchSecondGroup
	EventsRelaySwitchAllGroups
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
