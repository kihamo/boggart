package neptun

func (n *Neptun) SetModuleConfiguration(cfg *ModuleConfiguration) (err error) {
	_, err = n.client.WriteSingleRegister(AddressModuleConfiguration, cfg.Value())
	return err
}

func (n *Neptun) SetInputLinesConfiguration(l1, l2, l3, l4 *InputLinesConfiguration) (err error) {
	_, err = n.client.WriteSingleRegisterUint16Bytes(AddressInputLines12Configuration, []byte{l1.Value(), l2.Value()})
	if err != nil {
		return err
	}

	_, err = n.client.WriteSingleRegisterUint16Bytes(AddressInputLines34Configuration, []byte{l3.Value(), l4.Value()})
	return err
}

func (n *Neptun) SetEventsRelayConfiguration(cfg *EventsRelayConfiguration) (err error) {
	_, err = n.client.WriteSingleRegister(AddressEventsRelayConfiguration, uint16(cfg.Value()))
	return err
}

func (n *Neptun) SetCounterValue(counter, slot int, valueHigh, valueLow uint16) error {
	addressHigh, addressLow, err := n.counterValueAddresses(counter, slot)
	if err != nil {
		return err
	}

	_, err = n.client.WriteSingleRegister(addressHigh, valueHigh)
	if err != nil {
		return err
	}

	_, err = n.client.WriteSingleRegister(addressLow, valueLow)
	return err
}

func (n *Neptun) SetCounterConfiguration(cfg *CounterConfiguration, counter, slot int) error {
	address, err := n.counterConfigurationAddress(counter, slot)
	if err != nil {
		return err
	}

	_, err = n.client.WriteSingleRegister(address, cfg.Value())
	return err
}
