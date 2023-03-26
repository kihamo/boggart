package neptun

import (
	"encoding/binary"
	"fmt"
)

func (n *Neptun) ModuleConfiguration() (*ModuleConfiguration, error) {
	response, err := n.client.Read(AddressModuleConfiguration, 1)

	if err != nil {
		return nil, err
	}

	fmt.Println(response)

	return &ModuleConfiguration{
		value: binary.BigEndian.Uint16(response),
	}, err
}

func (n *Neptun) InputLines12Configuration() (err error) {
	response, err := n.client.Read(AddressInputLines12Configuration, 1)

	if err != nil {
		return err
	}

	fmt.Println(response)

	return err
}

func (n *Neptun) InputLinesStatus() (l1, l2, l3, l4 bool, err error) {
	response, err := n.client.Read(AddressInputLinesStatus, 1)

	if err == nil {
		l1 = response[0] != 0 // TODO:

		return l1, l2, l3, l4, err
	}

	return l1, l2, l3, l4, err
}

func (n *Neptun) EventsRelayConfiguration() (close, alarm uint8, err error) {
	response, err := n.client.Read(AddressEventsRelayConfiguration, 1)

	if err == nil {
		close = response[0]
		alarm = response[1]
	}

	return close, alarm, err
}

func (n *Neptun) SlaveIDAndBaudRate() (slaveId uint8, baudRate int, err error) {
	response, err := n.client.Read(AddressSlaveIDAndBaudRate, 1)

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

func (n *Neptun) Counter1Value() error {
	high, err := n.client.Read(AddressCounter1ValueHigh, 1)
	if err != nil {
		return err
	}

	low, err := n.client.Read(AddressCounter1ValueHigh, 1)
	if err != nil {
		return err
	}

	fmt.Println(high, low)

	return nil
}
