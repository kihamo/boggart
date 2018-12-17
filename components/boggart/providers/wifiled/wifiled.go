package wifiled

const (
	PortControlLocal  = 5577
	PortControlRemote = 80
	PortDiscover      = 48899

	// TODO: 0x21
	// TODO: 0x22
	// TODO: 0x31
	// TODO: 0x41
	// TODO: 0x51
	// TODO: 0x52

	CommandRemote   = 0xf0
	CommandLocal    = 0x0f
	CommandTimeSet  = 0x10
	CommandTimeSet2 = 0x14
	CommandTimeGet  = 0x11
	CommandTimeGet2 = 0x1a
	CommandTimeGet3 = 0x1b
	CommandMode     = 0x61
	CommandPower    = 0x71
	CommandState    = 0x81
	CommandState2   = 0x8a
	CommandState3   = 0x8b

	SpeedFast = 0x01 // 1
	SpeedSlow = 0x1f // 31
)

const (
	ModePreset1 Mode = 0x25 + iota
	ModePreset2
	ModePreset3
	ModePreset4
	ModePreset5
	ModePreset6
	ModePreset7
	ModePreset8
	ModePreset9
	ModePreset10
	ModePreset11
	ModePreset12
	ModePreset13
	ModePreset14
	ModePreset15
	ModePreset16
	ModePreset17
	ModePreset18
	ModePreset19
	ModePreset20
	ModePreset21 // not documented (smooth button)
)

const (
	ModeCustom Mode = 0x60 + iota
	ModeStatic
	ModeMusic
	ModeTesting
)

const (
	PowerOn uint8 = 0x23 + iota
	PowerOff
)

type Mode uint8

func (m Mode) String() string {
	switch m {
	case ModePreset1:
		return "7 colors gradual change"
	case ModePreset2:
		return "red gradual change"
	case ModePreset3:
		return "green gradual change"
	case ModePreset4:
		return "blue gradual change"
	case ModePreset5:
		return "yellow gradual change"
	case ModePreset6:
		return "cyan gradual change"
	case ModePreset7:
		return "purple gradual change"
	case ModePreset8:
		return "white gradual change"
	case ModePreset9:
		return "red and green gradual change"
	case ModePreset10:
		return "red and blue gradual change"
	case ModePreset11:
		return "green and blue gradual change"
	case ModePreset12:
		return "7 colors stroboflash"
	case ModePreset13:
		return "red stroboflash"
	case ModePreset14:
		return "green stroboflash"
	case ModePreset15:
		return "blue stroboflash"
	case ModePreset16:
		return "yellow stroboflash"
	case ModePreset17:
		return "cyan stroboflash"
	case ModePreset18:
		return "purpure stroboflash"
	case ModePreset19:
		return "white stroboflash"
	case ModePreset20:
		return "7 colors jump change"
	case ModePreset21:
		return "red, green, blue smooth change"
	case ModeCustom:
		return "custom"
	case ModeStatic:
		return "static"
	case ModeMusic:
		return "music"
	case ModeTesting:
		return "testing"
	default:
		return "unknown"
	}
}

type State struct {
	DeviceName uint8
	Power      bool
	Mode       Mode
	Speed      uint8
	Color      Color
}

type Color struct {
	Red          uint8
	Green        uint8
	Blue         uint8
	WarmWhite    uint8
	UseRGB       bool
	UseWarmWhite bool
}
