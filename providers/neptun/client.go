package neptun

import (
	"errors"
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

func (n *Neptun) counterValueAddresses(counter, slot int) (addressHigh, addressLow uint16, _ error) {
	if counter < 1 || counter > 2 {
		return 0, 0, errors.New("wrong counter number, only between 1 and 2")
	}

	if slot < 1 || slot > 4 {
		return 0, 0, errors.New("wrong slot number, only between 1 and 4")
	}

	switch slot {
	case 1:
		if counter == 1 {
			addressHigh = Counter1Slot1HighValue
			addressLow = Counter1Slot1LowValue
		} else {
			addressHigh = Counter2Slot1HighValue
			addressLow = Counter2Slot1LowValue
		}

	case 2:
		if counter == 1 {
			addressHigh = Counter1Slot2HighValue
			addressLow = Counter1Slot2LowValue
		} else {
			addressHigh = Counter2Slot2HighValue
			addressLow = Counter2Slot2LowValue
		}

	case 3:
		if counter == 1 {
			addressHigh = Counter1Slot3HighValue
			addressLow = Counter1Slot3LowValue
		} else {
			addressHigh = Counter2Slot3HighValue
			addressLow = Counter2Slot3LowValue
		}

	case 4:
		if counter == 1 {
			addressHigh = Counter1Slot4HighValue
			addressLow = Counter1Slot4LowValue
		} else {
			addressHigh = Counter2Slot4HighValue
			addressLow = Counter2Slot4LowValue
		}
	}

	return addressHigh, addressLow, nil
}

func (n *Neptun) counterConfigurationAddress(counter, slot int) (address uint16, _ error) {
	if counter < 1 || counter > 2 {
		return 0, errors.New("wrong counter number, only between 1 and 2")
	}

	if slot < 1 || slot > 4 {
		return 0, errors.New("wrong slot number, only between 1 and 4")
	}

	switch slot {
	case 1:
		if counter == 1 {
			address = Counter1Slot1Configuration
		} else {
			address = Counter2Slot1Configuration
		}

	case 2:
		if counter == 1 {
			address = Counter1Slot2Configuration
		} else {
			address = Counter2Slot2Configuration
		}

	case 3:
		if counter == 1 {
			address = Counter1Slot3Configuration
		} else {
			address = Counter2Slot3Configuration
		}

	case 4:
		if counter == 1 {
			address = Counter1Slot4Configuration
		} else {
			address = Counter2Slot4Configuration
		}
	}

	return address, nil
}
