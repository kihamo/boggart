package neptun

import (
	"encoding/binary"
	"fmt"
)

func (n *Neptun) ModuleConfiguration() (*ModuleConfiguration, error) {
	value, err := n.client.ReadHoldingRegistersUint16(AddressModuleConfiguration)

	if err != nil {
		return nil, err
	}

	return NewModuleConfiguration(uint(value)), err
}

func (n *Neptun) InputLinesConfiguration() (l1, l2, l3, l4 *InputLinesConfiguration, err error) {
	response, err := n.client.ReadHoldingRegisters(AddressInputLines12Configuration, 1)
	if err == nil {
		l1 = NewInputLinesConfiguration(uint(response[0]))
		l2 = NewInputLinesConfiguration(uint(response[1]))
	}

	response, err = n.client.ReadHoldingRegisters(AddressInputLines34Configuration, 1)
	if err == nil {
		l3 = NewInputLinesConfiguration(uint(response[0]))
		l4 = NewInputLinesConfiguration(uint(response[1]))
	}

	return l1, l2, l3, l4, err
}

func (n *Neptun) InputLinesStatus() (l1, l2, l3, l4 bool, err error) {
	value, err := n.client.ReadHoldingRegistersUint16(AddressInputLinesStatus)

	if err == nil {
		l1 = value&0 != 0
		l2 = value&1 != 0
		l3 = value&3 != 0
		l4 = value&4 != 0
	}

	return l1, l2, l3, l4, err
}

func (n *Neptun) EventsRelayConfiguration() (*EventsRelayConfiguration, error) {
	value, err := n.client.ReadHoldingRegistersUint16(AddressEventsRelayConfiguration)

	if err != nil {
		return nil, err
	}

	return NewEventsRelayConfiguration(uint(value)), err
}

func (n *Neptun) SlaveIDAndBaudRate() (slaveId uint8, baudRate int, err error) {
	response, err := n.client.ReadHoldingRegisters(AddressSlaveIDAndBaudRate, 1)

	if err != nil {
		return slaveId, baudRate, err
	}

	slaveId = response[0]

	switch response[1] {
	case 0x00:
		baudRate = 1200
	case 0x01:
		baudRate = 2400
	case 0x02:
		baudRate = 4800
	case 0x03:
		baudRate = 9600
	case 0x04:
		baudRate = 19200
	case 0x05:
		baudRate = 38400
	case 0x06:
		baudRate = 57600
	case 0x07:
		baudRate = 115200
	case 0x08:
		baudRate = 230400
	case 0x09:
		baudRate = 460800
	case 0x0A:
		baudRate = 921600
	}

	return slaveId, baudRate, err
}

func (n *Neptun) WirelessSensorCount() (uint16, error) {
	return n.client.ReadHoldingRegistersUint16(AddressWirelessSensorCount)
}

func (n *Neptun) WirelessSensorConfiguration(number int) (uint, error) {
	address, err := n.wirelessSensorConfigurationAddress(number)
	if err != nil {
		return 0, err
	}

	value, err := n.client.ReadHoldingRegistersUint8(address)
	if err != nil {
		return 0, err
	}

	return uint(value), err
}

func (n *Neptun) WirelessSensorStatus(number int) (*WirelessSensorStatus, error) {
	address, err := n.wirelessSensorStatusAddress(number)
	if err != nil {
		return nil, err
	}

	value, err := n.client.ReadHoldingRegistersUint16(address)
	if err != nil {
		return nil, err
	}

	return NewWirelessSensorStatus(uint(value)), err
}

func (n *Neptun) CountersValues() ([]*CounterValue, error) {
	const quantity uint16 = 16

	// counters(2) * slots(4) * bytes(2) = 16
	response, err := n.client.ReadHoldingRegisters(AddressCounter1Slot1HighValue, quantity)
	if err != nil {
		return nil, err
	}

	if len(response) != int(quantity)*2 {
		return nil, fmt.Errorf("wrong response payload length %d need %d", len(response), quantity*2)
	}

	ret := make([]*CounterValue, 0, 8)
	counter := 1
	slot := 1

	// counter 1 slot 1 - 0107 + 0108 00-01
	// counter 2 slot 1 - 0109 + 0110 02-03
	// counter 1 slot 2 - 0111 + 0112 04-05
	// counter 2 slot 2 - 0113 + 0114 06-07
	// counter 1 slot 3 - 0115 + 0116 08-09
	// counter 2 slot 3 - 0117 + 0118 10-11
	// counter 1 slot 4 - 0119 + 0120 12-13
	// counter 2 slot 4 - 0121 + 0122 14-15
	for i := 0; i < len(response); i += 4 {
		ret = append(ret, &CounterValue{
			number: counter,
			slot:   slot,
			value:  float64(binary.BigEndian.Uint32(response[i:i+4])) / 1000,
		})

		if counter != 2 {
			counter++
		} else {
			counter = 1

			if slot != 4 {
				slot++
			} else {
				slot = 1
			}
		}
	}

	return ret, nil
}

func (n *Neptun) CountersConfigurations() ([]*CounterConfiguration, error) {
	const quantity uint16 = 8

	// counters(2) * slots(4) = 8
	response, err := n.client.ReadHoldingRegisters(AddressCounter1Slot1Configuration, quantity)
	if err != nil {
		return nil, err
	}

	if len(response) != int(quantity)*2 {
		return nil, fmt.Errorf("wrong response payload length %d need %d", len(response), int(quantity)*2)
	}

	ret := make([]*CounterConfiguration, 0, 8)
	counter := 1
	slot := 1

	for i := 0; i < len(response); i += 2 {
		ret = append(ret, newCounterConfiguration(counter, slot, uint(binary.BigEndian.Uint16(response[i:i+2]))))

		if counter != 2 {
			counter++
		} else {
			counter = 1

			if slot != 4 {
				slot++
			} else {
				slot = 1
			}
		}
	}

	return ret, nil
}
