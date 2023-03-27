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
	Counter1Slot1HighValue           uint16 = 107
	Counter1Slot1LowValue            uint16 = 108
	Counter2Slot1HighValue           uint16 = 109
	Counter2Slot1LowValue            uint16 = 110
	Counter1Slot2HighValue           uint16 = 111
	Counter1Slot2LowValue            uint16 = 112
	Counter2Slot2HighValue           uint16 = 113
	Counter2Slot2LowValue            uint16 = 114
	Counter1Slot3HighValue           uint16 = 115
	Counter1Slot3LowValue            uint16 = 116
	Counter2Slot3HighValue           uint16 = 117
	Counter2Slot3LowValue            uint16 = 118
	Counter1Slot4HighValue           uint16 = 119
	Counter1Slot4LowValue            uint16 = 120
	Counter2Slot4HighValue           uint16 = 121
	Counter2Slot4LowValue            uint16 = 122
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

func (n *Neptun) Close() error {
	return n.client.Close()
}
