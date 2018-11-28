package broadlink

import (
	"encoding/binary"
	"errors"
	"net"

	"github.com/kihamo/boggart/components/boggart/providers/broadlink/internal"
)

type RM3Mini struct {
	*internal.Device
}

func NewRM3Mini(mac net.HardwareAddr, addr, iface net.UDPAddr) *RM3Mini {
	return &RM3Mini{
		Device: internal.NewDevice(KindRM3Mini, mac, addr, iface),
	}
}

func (d *RM3Mini) StartCaptureRemoteControlCode() error {
	payload := make([]byte, 0x10)
	payload[0] = 0x03

	response, err := d.Call(0x6a, payload)
	if err != nil {
		return err
	}

	code := binary.LittleEndian.Uint16(response[0x22:0x24])
	if code != 0 {
		return errors.New("failed to start capturing remote control code")
	}

	return nil
}

func (d *RM3Mini) ReadCapturedRemoteControlCode() ([]byte, error) {
	payload := make([]byte, 0x10)
	payload[0] = 4

	response, err := d.Call(0x6a, payload)
	if err != nil {
		return nil, err
	}

	code := binary.LittleEndian.Uint16(response[0x22:0x24])
	if code != 0 {
		if code == 0xfff6 {
			return nil, errors.New("signal not captured")
		}

		return nil, errors.New("failed to read capturing remote control code")
	}

	data, err := d.DecodePacket(response)
	if err != nil {
		return nil, err
	}

	if len(data) < 8 {
		return nil, errors.New("incomplete data")
	}

	cmd := binary.LittleEndian.Uint16(data[:4])
	if cmd != 0x04 {
		return nil, errors.New("invalid command code")
	}

	return data[0x04:], nil
}
