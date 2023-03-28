package neptun

import (
	"errors"
)

func (n *Neptun) ModuleConfiguration() (*ModuleConfiguration, error) {
	value, err := n.client.ReadHoldingRegistersUint16(AddressModuleConfiguration)

	if err != nil {
		return nil, err
	}

	return &ModuleConfiguration{
		value: value,
	}, err
}

func (n *Neptun) InputLinesConfiguration() (l1, l2, l3, l4 *InputLinesConfiguration, err error) {
	response, err := n.client.ReadHoldingRegisters(AddressInputLines12Configuration, 1)
	if err == nil {
		l1 = &InputLinesConfiguration{
			value: response[0],
		}
		l2 = &InputLinesConfiguration{
			value: response[1],
		}
	}

	response, err = n.client.ReadHoldingRegisters(AddressInputLines34Configuration, 1)
	if err == nil {
		l3 = &InputLinesConfiguration{
			value: response[0],
		}
		l4 = &InputLinesConfiguration{
			value: response[1],
		}
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

	return &EventsRelayConfiguration{
		value: value,
	}, err
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

func (n *Neptun) CounterValue(counter, slot int) (uint16, uint16, error) {
	addressHigh, addressLow, err := n.counterAddresses(counter, slot)
	if err != nil {
		return 0, 0, err
	}

	valueHigh, err := n.client.ReadHoldingRegistersUint16(addressHigh)
	if err != nil {
		return 0, 0, err
	}

	valueLow, err := n.client.ReadHoldingRegistersUint16(addressLow)
	if err != nil {
		return 0, 0, err
	}

	return valueHigh, valueLow, nil
}

func (n *Neptun) CounterConfiguration(counter, slot int) (*CounterConfiguration, error) {
	if counter < 1 || counter > 2 {
		return nil, errors.New("wrong counter number, only between 1 and 2")
	}

	if slot < 1 || slot > 4 {
		return nil, errors.New("wrong slot number, only between 1 and 4")
	}

	var address uint16

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

	value, err := n.client.ReadHoldingRegistersUint16(address)

	if err != nil {
		return nil, err
	}

	return &CounterConfiguration{
		value: value,
	}, err
}
