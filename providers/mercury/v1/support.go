package v1

import (
	"errors"
)

const (
	DeviceUnknown = uint8(0) + iota
	Device200
	Device201
	Device203
	Device206
)

var (
	supported map[uint8]map[uint8]bool

	ErrCommandNotSupported = errors.New("command not supported")
)

func init() {
	supported = map[uint8]map[uint8]bool{
		Device200: {
			RequestCommandReadParamLastChange: false,
			RequestCommandReadMaximum:         false,
			RequestCommandReadDisplayModeExt:  false,
		},
	}
}

func CommandNotSupported(err error) bool {
	return err == ErrCommandNotSupported
}

func IsCommandSupported(device, command uint8) bool {
	if dev, ok := supported[device]; ok {
		if cmd, ok := dev[command]; ok {
			return cmd
		}
	}

	return true
}
