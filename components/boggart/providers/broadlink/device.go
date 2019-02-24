package broadlink

import (
	"net"
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/broadlink/internal"
)

type RemoteType int

/*
{ "type": "0x2711", "name": "SP2"},
{ "type": "0x753e", "name": "SP3"},
{ "type": "0x947a", "name": "SP3S_EU"},
{ "type": "0x9479", "name": "SP3S_US"},
{ "type": "0x2720", "name": "SP_MINI"},
{ "type": "0x7549", "name": "SP_MINI2_CHINA_MOBILE"},
{ "type": "0x7918", "name": "SP_MINI2_OEM_MAX"},
{ "type": "0x7530", "name": "SP_MINI2_OEM_MIN"},
{ "type": "0x7547", "name": "SP_MINI2_WIFI_BOX"},
{ "type": "0x2733", "name": "SP_MINI_CC"},
{ "type": "0x7539", "name": "SP_MINI_HAIBEI"},
{ "type": "0x754e", "name": "SP_MINI_HYC"},
{ "type": "0x753d", "name": "SP_MINI_KPL"},
{ "type": "0x7536", "name": "SP_MINI_NEO"},
{ "type": "0x273e", "name": "SP_MINI_PHICOMM"},
{ "type": "0x2736", "name": "SP_MINI_PLUS"},
{ "type": "0x947c", "name": "SP_MINI_PLUS2"},
{ "type": "0x2728", "name": "SP_MINI_V2"}

0: SP1
0x2711: SP2
0x2719, 0x7919, 0x271a, 0x791a: Honeywell SP2
0x2720: SPMini
0x753e: SP3
0x2728: SPMini2
0x2733, 0x273e, 0x7539, 0x754e, 0x753d, 0x7536: OEM branded SPMini
0x7540: MP2
0x7530, 0x7918, 0x7549: OEM branded SPMini2
0x2736: SPMiniPlus
0x947c: SPMiniPlus
0x7547: SC1
0x947a, 0x9479: SP3S
0x2710: RM1
0x2712: RM2
0x2737: RM Mini
0x27a2: RM Mini R2
0x273d: RM Pro Phicomm
0x2783: RM2 Home Plus
0x277c: RM2 Home Plus GDT
0x272a: RM2 Pro Plus
0x2787: RM2 Pro Plus2
0x279d: RM2 Pro Plus3
0x2797: RM2 Pro Plus HYC
0x278b: RM2 Pro Plus BL
0x27a1: RM2 Pro Plus R1
0x278f: RM Mini Shate
0x2714, 0x27a3: A1
0x4EB5: MP1
0x271F: MS1
0x2722: S1
0x273c: S1 Phicomm
0x4f34, 0x4f35, 0x4f36: TW2 Switch
0x4ee6, 0x4eee, 0x4eef: NEW Switch
0x271b, 0x271c: Honyar switch
0x2721: Camera
0x42, 0x4e62: DEYE HUMIDIFIER
0x2d, 0x4f42, 0x4e4d: DOOYA CURTAIN
0x2723, 0x4eda: HONYAR MS
0x2727, 0x2726, 0x2724, 0x2725: HONYAR SL
0x4c, 0x4e6c: MFRESH AIR
0x271e, 0x2746: PLC (TW_ROUTER)
0x2774, 0x7530, 0x2742, 0x4e20: MIN/MAX AP/OEM
0x4e69: LIGHTMATES
*/

const (
	KindRMMini      = 0x2737
	KindRM2ProPlus3 = 0x279d
	KindSP3SEU      = 0x947a
	KindSP3SUS      = 0x9479

	RemoteIR       RemoteType = 0x26
	RemoteRF433Mhz RemoteType = 0xb2
	RemoteRF315Mhz RemoteType = 0xd7
)

type Device interface {
	Kind() int
	MAC() net.HardwareAddr
	Addr() *net.UDPAddr
	Interface() *net.UDPAddr
	SetTimeout(duration time.Duration)
}

func NewDevice(kind int, mac net.HardwareAddr, addr, iface net.UDPAddr) Device {
	switch kind {
	case KindSP3SEU:
		return NewSP3SEU(mac, addr, iface)

	case KindSP3SUS:
		return NewSP3SUS(mac, addr, iface)

	case KindRM2ProPlus3:
		return NewRM2ProPlus3(mac, addr, iface)

	case KindRMMini:
		return NewRMMini(mac, addr, iface)
	}

	return internal.NewDevice(kind, mac, addr, iface)
}
