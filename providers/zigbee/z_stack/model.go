package z_stack

import (
	"encoding/hex"
	"sync"
)

type Device struct {
	networkAddress uint16
	ieeeAddress    []byte
	capabilities   uint8

	lock sync.RWMutex
}

func (d *Device) NetworkAddress() uint16 {
	d.lock.RLock()
	defer d.lock.RUnlock()

	return d.networkAddress
}

func (d *Device) SetNetworkAddress(address uint16) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.networkAddress = address
}

func (d *Device) IEEEAddress() []byte {
	d.lock.RLock()
	defer d.lock.RUnlock()

	return d.ieeeAddress
}

func (d *Device) SetIEEEAddress(address []byte) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.ieeeAddress = address
}

func (d *Device) IEEEAddressAsString() string {
	return hex.EncodeToString(d.IEEEAddress())
}

func (d *Device) Capabilities() uint8 {
	d.lock.RLock()
	defer d.lock.RUnlock()

	return d.capabilities
}

func (d *Device) SetCapabilities(capabilities uint8) {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.capabilities = capabilities
}
