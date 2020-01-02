// +build !linux,!darwin

package bluetooth

import (
	"errors"

	"github.com/go-ble/ble"
)

// NewDevice ...
func NewDevice(opts ...ble.Option) (d ble.Device, err error) {
	return nil, errors.New("not implemented")
}
