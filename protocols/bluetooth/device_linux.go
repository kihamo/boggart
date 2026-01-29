//go:build linux
// +build linux

package bluetooth

import (
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

// DefaultDevice ...
func NewDevice(opts ...ble.Option) (d ble.Device, err error) {
	return linux.NewDevice(opts...)
}
