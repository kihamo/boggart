package broadlink

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/kihamo/boggart/components/boggart/providers/broadlink/internal"
)

type SP3S struct {
	*internal.Device
}

func NewSP3S(mac net.HardwareAddr, addr, iface net.UDPAddr) *SP3S {
	return &SP3S{
		Device: internal.NewDevice(KindSP3S, mac, addr, iface),
	}
}

func (d *SP3S) PowerOn() error {
	payload := make([]byte, 16)
	payload[0] = 2
	payload[4] = 1

	return d.Cmd(0x6a, payload)
}

func (d *SP3S) PowerOff() error {
	payload := make([]byte, 16)
	payload[0] = 2
	payload[4] = 0

	return d.Cmd(0x6a, payload)
}

func (d *SP3S) PowerState() (bool, error) {
	payload := make([]byte, 16)
	payload[0] = 1

	response, err := d.Call(0x6a, payload)
	if err != nil {
		return false, err
	}

	data, err := d.DecodePacket(response)
	if err != nil {
		return false, err
	}

	return data[4] != 0, nil
}

func (d *SP3S) Energy() (float64, error) {
	payload := []byte{8, 0, 254, 1, 5, 1, 0, 0, 0, 45}
	response, err := d.Call(0x6a, payload)
	if err != nil {
		return -1, err
	}

	rescode := binary.LittleEndian.Uint16(response[0x22:0x24])
	if rescode != 0 {
		return -1, errors.New("response code isn't 0")
	}

	data, err := d.DecodePacket(response)
	if err != nil {
		return -1, err
	}

	fmt.Printf("%X\n", data[7])
	fmt.Printf("%X\n", data[6])
	fmt.Printf("%X\n", data[5])

	// int(hex(payload[0x07] * 256 + payload[0x06])[2:]) + int(hex(payload[0x05])[2:])/100.0

	return -1, nil
}
