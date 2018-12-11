package broadlink

import (
	"encoding/hex"
	"net"

	"github.com/kihamo/boggart/components/boggart/providers/broadlink/internal"
)

type RMProPlus struct {
	*RM3Mini
}

func NewRMProPlus(mac net.HardwareAddr, addr, iface net.UDPAddr) *RMProPlus {
	return &RMProPlus{
		RM3Mini: &RM3Mini{
			Device: internal.NewDevice(KindRMProPlus, mac, addr, iface),
		},
	}
}

func (d *RMProPlus) SendRF433MhzRemoteControlCode(code []byte, count int) error {
	return d.SendRemoteControlCode(RemoteRF433Mhz, code, count)
}

func (d *RMProPlus) SendRF433MhzRemoteControlCodeAsString(code string, count int) error {
	decoded, err := hex.DecodeString(code)
	if err != nil {
		return err
	}

	return d.SendRemoteControlCode(RemoteRF433Mhz, decoded, count)
}

func (d *RMProPlus) SendRF315MhzRemoteControlCode(code []byte, count int) error {
	return d.SendRemoteControlCode(RemoteRF315Mhz, code, count)
}

func (d *RMProPlus) SendRF315MhzRemoteControlCodeAsString(code string, count int) error {
	decoded, err := hex.DecodeString(code)
	if err != nil {
		return err
	}

	return d.SendRemoteControlCode(RemoteRF315Mhz, decoded, count)
}
