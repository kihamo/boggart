package neptun

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
