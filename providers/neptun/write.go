package neptun

func (n *Neptun) SetModuleConfiguration(cfg *ModuleConfiguration) (err error) {
	_, err = n.client.WriteSingleRegister(AddressModuleConfiguration, cfg.value)
	return err
}

func (n *Neptun) SetEventsRelayConfiguration(cfg *EventsRelayConfiguration) (err error) {
	_, err = n.client.WriteSingleRegister(AddressEventsRelayConfiguration, cfg.value)
	return err
}
