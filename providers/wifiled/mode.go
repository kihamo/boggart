package wifiled

import "strconv"

type Mode uint8

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

func ModeFromString(mode string) (*Mode, error) {
	v, err := strconv.ParseUint(mode, 10, 64)
	if err != nil {
		return nil, err
	}

	m := Mode(v)
	return &m, nil
}
