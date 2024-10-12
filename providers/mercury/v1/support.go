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

/*
200.00 -- однотарифный, без доп функций
200.02 -- многотарифный, интерфейс CAN (установлен этот)
200.03 -- многотарифный, отключение нагрузки, интерфейс CAN
200.04 -- многотарифный, отключение нагрузки, интерфейс CAN, модем PLC
*/

func init() {
	supported = map[uint8]map[uint8]bool{
		Device200: {
			CommandReadParamLastChange: false,
			CommandReadRelayMode:       false,
			CommandReadMaximum:         false,
			CommandReadDisplayModeExt:  false,
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
