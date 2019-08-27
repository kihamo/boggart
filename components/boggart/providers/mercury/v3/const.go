package v3

const (
	AddressUniveral      = 0x00
	AddressPersonalBegin = 0x01
	AddressPersonalEnd   = 0xF0
	AddressBroadCast     = 0xFE
	AddressReservedBegin = 0xF1
	AddressReservedEnd   = 0xFF
)

type LevelPassword [6]byte

var (
	DefaultPasswordLevel1 = LevelPassword{0x1, 0x1, 0x1, 0x1, 0x1, 0x1}
	DefaultPasswordLevel2 = LevelPassword{0x2, 0x2, 0x2, 0x2, 0x2, 0x2}
)

func (p LevelPassword) Bytes() []byte {
	return p[:]
}
