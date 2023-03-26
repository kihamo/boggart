package neptun

import (
	"encoding/binary"
)

func (n *Neptun) SetModuleConfiguration(cfg *ModuleConfiguration) (err error) {
	_, err = n.client.WriteSingleRegister(AddressModuleConfiguration, cfg.value)
	return err
}

func (n *Neptun) SetEventsRelayConfiguration(close, alarm uint8) (err error) {
	_, err = n.client.WriteSingleRegister(AddressEventsRelayConfiguration, binary.BigEndian.Uint16([]byte{close, alarm}))
	return err
}
