package wifiled

const (
	PortControlLocal  = 5577
	PortControlRemote = 80
	PortDiscover      = 48899

	True  = 0xf0
	False = 0x0f

	// TODO: 0x21
	// TODO: 0x22
	// TODO: 0x51
	// TODO: 0x52

	CommandRemote = True
	CommandLocal  = False

	CommandTimeSet         = 0x10
	CommandTimeSet2        = 0x14
	CommandTimeGet         = 0x11
	CommandTimeGet2        = 0x1a
	CommandTimeGet3        = 0x1b
	CommandColorSetPersist = 0x31
	CommandColorSet        = 0x41
	CommandMode            = 0x61
	CommandPower           = 0x71
	CommandState           = 0x81
	CommandState2          = 0x8a
	CommandState3          = 0x8b

	SpeedFast = 0x01 // 1
	SpeedSlow = 0x1f // 31
)

const (
	PowerOn uint8 = 0x23 + iota
	PowerOff
)

type State struct {
	DeviceName uint8
	Power      bool
	Mode       Mode
	Speed      uint8
	Color      Color
}
