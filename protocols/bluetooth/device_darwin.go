// +build darwin

package bluetooth

import (
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/darwin"
)

// DefaultDevice ...
func NewDevice(opts ...ble.Option) (d ble.Device, err error) {
	return darwin.NewDevice(opts...)
}
