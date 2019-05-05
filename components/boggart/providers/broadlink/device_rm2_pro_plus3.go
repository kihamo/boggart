package broadlink

import (
	"encoding/hex"
	"net"

	"github.com/kihamo/boggart/components/boggart/providers/broadlink/internal"
)

type RM2ProPlus3 struct {
	*RMMini
}

func NewRM2ProPlus3(mac net.HardwareAddr, addr string) *RM2ProPlus3 {
	return &RM2ProPlus3{
		RMMini: &RMMini{
			Device: internal.NewDevice(KindRM2ProPlus3, mac, addr),
		},
	}
}

func (d *RM2ProPlus3) SendRF433MhzRemoteControlCode(code []byte, count int) error {
	return d.SendRemoteControlCode(RemoteRF433Mhz, code, count)
}

func (d *RM2ProPlus3) SendRF433MhzRemoteControlCodeAsString(code string, count int) error {
	decoded, err := hex.DecodeString(code)
	if err != nil {
		return err
	}

	return d.SendRemoteControlCode(RemoteRF433Mhz, decoded, count)
}
