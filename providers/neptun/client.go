package neptun

import (
	"errors"
	"net/url"

	"github.com/kihamo/boggart/protocols/modbus"
)

const (
	AddressModuleConfiguration           uint16 = 0
	AddressInputLines12Configuration     uint16 = 1
	AddressInputLines34Configuration     uint16 = 2
	AddressInputLinesStatus              uint16 = 3
	AddressEventsRelayConfiguration      uint16 = 4
	AddressSlaveIDAndBaudRate            uint16 = 5
	AddressWirelessSensorCount           uint16 = 6
	AddressWirelessSensor1Configuration  uint16 = 7
	AddressWirelessSensor2Configuration  uint16 = 8
	AddressWirelessSensor3Configuration  uint16 = 9
	AddressWirelessSensor4Configuration  uint16 = 10
	AddressWirelessSensor5Configuration  uint16 = 11
	AddressWirelessSensor6Configuration  uint16 = 12
	AddressWirelessSensor7Configuration  uint16 = 13
	AddressWirelessSensor8Configuration  uint16 = 14
	AddressWirelessSensor9Configuration  uint16 = 15
	AddressWirelessSensor10Configuration uint16 = 16
	AddressWirelessSensor11Configuration uint16 = 17
	AddressWirelessSensor12Configuration uint16 = 18
	AddressWirelessSensor13Configuration uint16 = 19
	AddressWirelessSensor14Configuration uint16 = 20
	AddressWirelessSensor15Configuration uint16 = 21
	AddressWirelessSensor16Configuration uint16 = 22
	AddressWirelessSensor17Configuration uint16 = 23
	AddressWirelessSensor18Configuration uint16 = 24
	AddressWirelessSensor19Configuration uint16 = 25
	AddressWirelessSensor20Configuration uint16 = 26
	AddressWirelessSensor21Configuration uint16 = 27
	AddressWirelessSensor22Configuration uint16 = 28
	AddressWirelessSensor23Configuration uint16 = 29
	AddressWirelessSensor24Configuration uint16 = 30
	AddressWirelessSensor25Configuration uint16 = 31
	AddressWirelessSensor26Configuration uint16 = 32
	AddressWirelessSensor27Configuration uint16 = 33
	AddressWirelessSensor28Configuration uint16 = 34
	AddressWirelessSensor29Configuration uint16 = 35
	AddressWirelessSensor30Configuration uint16 = 36
	AddressWirelessSensor31Configuration uint16 = 37
	AddressWirelessSensor32Configuration uint16 = 38
	AddressWirelessSensor33Configuration uint16 = 39
	AddressWirelessSensor34Configuration uint16 = 40
	AddressWirelessSensor35Configuration uint16 = 41
	AddressWirelessSensor36Configuration uint16 = 42
	AddressWirelessSensor37Configuration uint16 = 43
	AddressWirelessSensor38Configuration uint16 = 44
	AddressWirelessSensor39Configuration uint16 = 45
	AddressWirelessSensor40Configuration uint16 = 46
	AddressWirelessSensor41Configuration uint16 = 47
	AddressWirelessSensor42Configuration uint16 = 48
	AddressWirelessSensor43Configuration uint16 = 49
	AddressWirelessSensor44Configuration uint16 = 50
	AddressWirelessSensor45Configuration uint16 = 51
	AddressWirelessSensor46Configuration uint16 = 52
	AddressWirelessSensor47Configuration uint16 = 53
	AddressWirelessSensor48Configuration uint16 = 54
	AddressWirelessSensor49Configuration uint16 = 55
	AddressWirelessSensor50Configuration uint16 = 56
	AddressWirelessSensor1Status         uint16 = 57
	AddressWirelessSensor2Status         uint16 = 58
	AddressWirelessSensor3Status         uint16 = 59
	AddressWirelessSensor4Status         uint16 = 60
	AddressWirelessSensor5Status         uint16 = 61
	AddressWirelessSensor6Status         uint16 = 62
	AddressWirelessSensor7Status         uint16 = 63
	AddressWirelessSensor8Status         uint16 = 64
	AddressWirelessSensor9Status         uint16 = 65
	AddressWirelessSensor10Status        uint16 = 66
	AddressWirelessSensor11Status        uint16 = 67
	AddressWirelessSensor12Status        uint16 = 68
	AddressWirelessSensor13Status        uint16 = 69
	AddressWirelessSensor14Status        uint16 = 70
	AddressWirelessSensor15Status        uint16 = 71
	AddressWirelessSensor16Status        uint16 = 72
	AddressWirelessSensor17Status        uint16 = 73
	AddressWirelessSensor18Status        uint16 = 74
	AddressWirelessSensor19Status        uint16 = 75
	AddressWirelessSensor20Status        uint16 = 76
	AddressWirelessSensor21Status        uint16 = 77
	AddressWirelessSensor22Status        uint16 = 78
	AddressWirelessSensor23Status        uint16 = 79
	AddressWirelessSensor24Status        uint16 = 80
	AddressWirelessSensor25Status        uint16 = 81
	AddressWirelessSensor26Status        uint16 = 82
	AddressWirelessSensor27Status        uint16 = 83
	AddressWirelessSensor28Status        uint16 = 84
	AddressWirelessSensor29Status        uint16 = 85
	AddressWirelessSensor30Status        uint16 = 86
	AddressWirelessSensor31Status        uint16 = 87
	AddressWirelessSensor32Status        uint16 = 88
	AddressWirelessSensor33Status        uint16 = 89
	AddressWirelessSensor34Status        uint16 = 90
	AddressWirelessSensor35Status        uint16 = 91
	AddressWirelessSensor36Status        uint16 = 92
	AddressWirelessSensor37Status        uint16 = 93
	AddressWirelessSensor38Status        uint16 = 94
	AddressWirelessSensor39Status        uint16 = 95
	AddressWirelessSensor40Status        uint16 = 96
	AddressWirelessSensor41Status        uint16 = 97
	AddressWirelessSensor42Status        uint16 = 98
	AddressWirelessSensor43Status        uint16 = 99
	AddressWirelessSensor44Status        uint16 = 100
	AddressWirelessSensor45Status        uint16 = 101
	AddressWirelessSensor46Status        uint16 = 102
	AddressWirelessSensor47Status        uint16 = 103
	AddressWirelessSensor48Status        uint16 = 104
	AddressWirelessSensor49Status        uint16 = 105
	AddressWirelessSensor50Status        uint16 = 106
	AddressCounter1Slot1HighValue        uint16 = 107
	AddressCounter1Slot1LowValue         uint16 = 108
	AddressCounter2Slot1HighValue        uint16 = 109
	AddressCounter2Slot1LowValue         uint16 = 110
	AddressCounter1Slot2HighValue        uint16 = 111
	AddressCounter1Slot2LowValue         uint16 = 112
	AddressCounter2Slot2HighValue        uint16 = 113
	AddressCounter2Slot2LowValue         uint16 = 114
	AddressCounter1Slot3HighValue        uint16 = 115
	AddressCounter1Slot3LowValue         uint16 = 116
	AddressCounter2Slot3HighValue        uint16 = 117
	AddressCounter2Slot3LowValue         uint16 = 118
	AddressCounter1Slot4HighValue        uint16 = 119
	AddressCounter1Slot4LowValue         uint16 = 120
	AddressCounter2Slot4HighValue        uint16 = 121
	AddressCounter2Slot4LowValue         uint16 = 122
	AddressCounter1Slot1Configuration    uint16 = 123
	AddressCounter2Slot1Configuration    uint16 = 124
	AddressCounter1Slot2Configuration    uint16 = 125
	AddressCounter2Slot2Configuration    uint16 = 126
	AddressCounter1Slot3Configuration    uint16 = 127
	AddressCounter2Slot3Configuration    uint16 = 128
	AddressCounter1Slot4Configuration    uint16 = 129
	AddressCounter2Slot4Configuration    uint16 = 130
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

func (n *Neptun) wirelessSensorConfigurationAddress(number int) (address uint16, _ error) {
	switch number {
	case 1:
		return AddressWirelessSensor1Configuration, nil
	case 2:
		return AddressWirelessSensor2Configuration, nil
	case 3:
		return AddressWirelessSensor3Configuration, nil
	case 4:
		return AddressWirelessSensor4Configuration, nil
	case 5:
		return AddressWirelessSensor5Configuration, nil
	case 6:
		return AddressWirelessSensor6Configuration, nil
	case 7:
		return AddressWirelessSensor7Configuration, nil
	case 8:
		return AddressWirelessSensor8Configuration, nil
	case 9:
		return AddressWirelessSensor9Configuration, nil
	case 10:
		return AddressWirelessSensor10Configuration, nil
	case 11:
		return AddressWirelessSensor11Configuration, nil
	case 12:
		return AddressWirelessSensor12Configuration, nil
	case 13:
		return AddressWirelessSensor13Configuration, nil
	case 14:
		return AddressWirelessSensor14Configuration, nil
	case 15:
		return AddressWirelessSensor15Configuration, nil
	case 16:
		return AddressWirelessSensor16Configuration, nil
	case 17:
		return AddressWirelessSensor17Configuration, nil
	case 18:
		return AddressWirelessSensor18Configuration, nil
	case 19:
		return AddressWirelessSensor19Configuration, nil
	case 20:
		return AddressWirelessSensor20Configuration, nil
	case 21:
		return AddressWirelessSensor21Configuration, nil
	case 22:
		return AddressWirelessSensor22Configuration, nil
	case 23:
		return AddressWirelessSensor23Configuration, nil
	case 24:
		return AddressWirelessSensor24Configuration, nil
	case 25:
		return AddressWirelessSensor25Configuration, nil
	case 26:
		return AddressWirelessSensor26Configuration, nil
	case 27:
		return AddressWirelessSensor27Configuration, nil
	case 28:
		return AddressWirelessSensor28Configuration, nil
	case 29:
		return AddressWirelessSensor29Configuration, nil
	case 30:
		return AddressWirelessSensor30Configuration, nil
	case 31:
		return AddressWirelessSensor31Configuration, nil
	case 32:
		return AddressWirelessSensor32Configuration, nil
	case 33:
		return AddressWirelessSensor33Configuration, nil
	case 34:
		return AddressWirelessSensor34Configuration, nil
	case 35:
		return AddressWirelessSensor35Configuration, nil
	case 36:
		return AddressWirelessSensor36Configuration, nil
	case 37:
		return AddressWirelessSensor37Configuration, nil
	case 38:
		return AddressWirelessSensor38Configuration, nil
	case 39:
		return AddressWirelessSensor39Configuration, nil
	case 40:
		return AddressWirelessSensor40Configuration, nil
	case 41:
		return AddressWirelessSensor41Configuration, nil
	case 42:
		return AddressWirelessSensor42Configuration, nil
	case 43:
		return AddressWirelessSensor43Configuration, nil
	case 44:
		return AddressWirelessSensor44Configuration, nil
	case 45:
		return AddressWirelessSensor45Configuration, nil
	case 46:
		return AddressWirelessSensor46Configuration, nil
	case 47:
		return AddressWirelessSensor47Configuration, nil
	case 48:
		return AddressWirelessSensor48Configuration, nil
	case 49:
		return AddressWirelessSensor49Configuration, nil
	case 50:
		return AddressWirelessSensor50Configuration, nil

	default:
		return 0, errors.New("wrong wireless sensor number, only between 1 and 50")
	}
}
func (n *Neptun) wirelessSensorStatusAddress(number int) (address uint16, _ error) {
	switch number {
	case 1:
		return AddressWirelessSensor1Status, nil
	case 2:
		return AddressWirelessSensor2Status, nil
	case 3:
		return AddressWirelessSensor3Status, nil
	case 4:
		return AddressWirelessSensor4Status, nil
	case 5:
		return AddressWirelessSensor5Status, nil
	case 6:
		return AddressWirelessSensor6Status, nil
	case 7:
		return AddressWirelessSensor7Status, nil
	case 8:
		return AddressWirelessSensor8Status, nil
	case 9:
		return AddressWirelessSensor9Status, nil
	case 10:
		return AddressWirelessSensor10Status, nil
	case 11:
		return AddressWirelessSensor11Status, nil
	case 12:
		return AddressWirelessSensor12Status, nil
	case 13:
		return AddressWirelessSensor13Status, nil
	case 14:
		return AddressWirelessSensor14Status, nil
	case 15:
		return AddressWirelessSensor15Status, nil
	case 16:
		return AddressWirelessSensor16Status, nil
	case 17:
		return AddressWirelessSensor17Status, nil
	case 18:
		return AddressWirelessSensor18Status, nil
	case 19:
		return AddressWirelessSensor19Status, nil
	case 20:
		return AddressWirelessSensor20Status, nil
	case 21:
		return AddressWirelessSensor21Status, nil
	case 22:
		return AddressWirelessSensor22Status, nil
	case 23:
		return AddressWirelessSensor23Status, nil
	case 24:
		return AddressWirelessSensor24Status, nil
	case 25:
		return AddressWirelessSensor25Status, nil
	case 26:
		return AddressWirelessSensor26Status, nil
	case 27:
		return AddressWirelessSensor27Status, nil
	case 28:
		return AddressWirelessSensor28Status, nil
	case 29:
		return AddressWirelessSensor29Status, nil
	case 30:
		return AddressWirelessSensor30Status, nil
	case 31:
		return AddressWirelessSensor31Status, nil
	case 32:
		return AddressWirelessSensor32Status, nil
	case 33:
		return AddressWirelessSensor33Status, nil
	case 34:
		return AddressWirelessSensor34Status, nil
	case 35:
		return AddressWirelessSensor35Status, nil
	case 36:
		return AddressWirelessSensor36Status, nil
	case 37:
		return AddressWirelessSensor37Status, nil
	case 38:
		return AddressWirelessSensor38Status, nil
	case 39:
		return AddressWirelessSensor39Status, nil
	case 40:
		return AddressWirelessSensor40Status, nil
	case 41:
		return AddressWirelessSensor41Status, nil
	case 42:
		return AddressWirelessSensor42Status, nil
	case 43:
		return AddressWirelessSensor43Status, nil
	case 44:
		return AddressWirelessSensor44Status, nil
	case 45:
		return AddressWirelessSensor45Status, nil
	case 46:
		return AddressWirelessSensor46Status, nil
	case 47:
		return AddressWirelessSensor47Status, nil
	case 48:
		return AddressWirelessSensor48Status, nil
	case 49:
		return AddressWirelessSensor49Status, nil
	case 50:
		return AddressWirelessSensor50Status, nil

	default:
		return 0, errors.New("wrong wireless sensor number, only between 1 and 50")
	}
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
			addressHigh = AddressCounter1Slot1HighValue
			addressLow = AddressCounter1Slot1LowValue
		} else {
			addressHigh = AddressCounter2Slot1HighValue
			addressLow = AddressCounter2Slot1LowValue
		}

	case 2:
		if counter == 1 {
			addressHigh = AddressCounter1Slot2HighValue
			addressLow = AddressCounter1Slot2LowValue
		} else {
			addressHigh = AddressCounter2Slot2HighValue
			addressLow = AddressCounter2Slot2LowValue
		}

	case 3:
		if counter == 1 {
			addressHigh = AddressCounter1Slot3HighValue
			addressLow = AddressCounter1Slot3LowValue
		} else {
			addressHigh = AddressCounter2Slot3HighValue
			addressLow = AddressCounter2Slot3LowValue
		}

	case 4:
		if counter == 1 {
			addressHigh = AddressCounter1Slot4HighValue
			addressLow = AddressCounter1Slot4LowValue
		} else {
			addressHigh = AddressCounter2Slot4HighValue
			addressLow = AddressCounter2Slot4LowValue
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
			address = AddressCounter1Slot1Configuration
		} else {
			address = AddressCounter2Slot1Configuration
		}

	case 2:
		if counter == 1 {
			address = AddressCounter1Slot2Configuration
		} else {
			address = AddressCounter2Slot2Configuration
		}

	case 3:
		if counter == 1 {
			address = AddressCounter1Slot3Configuration
		} else {
			address = AddressCounter2Slot3Configuration
		}

	case 4:
		if counter == 1 {
			address = AddressCounter1Slot4Configuration
		} else {
			address = AddressCounter2Slot4Configuration
		}
	}

	return address, nil
}
