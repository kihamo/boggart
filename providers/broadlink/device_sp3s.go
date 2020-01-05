package broadlink

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"github.com/kihamo/boggart/providers/broadlink/internal"
)

type SP3S struct {
	*internal.Device
}

func NewSP3SEU(mac net.HardwareAddr, addr string) *SP3S {
	return &SP3S{
		Device: internal.NewDevice(KindSP3SEU, mac, addr),
	}
}

func NewSP3SUS(mac net.HardwareAddr, addr string) *SP3S {
	return &SP3S{
		Device: internal.NewDevice(KindSP3SUS, mac, addr),
	}
}

func (d *SP3S) On() error {
	payload := make([]byte, 16)
	payload[0] = 2
	payload[4] = 1

	return d.Cmd(0x6a, payload)
}

func (d *SP3S) Off() error {
	payload := make([]byte, 16)
	payload[0] = 2
	payload[4] = 0

	return d.Cmd(0x6a, payload)
}

func (d *SP3S) State() (bool, error) {
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

func (d *SP3S) Power() (float64, error) {
	payload := []byte{8, 0, 254, 1, 5, 1, 0, 0, 0, 45}
	response, err := d.Call(0x6a, payload)
	if err != nil {
		return -1, err
	}

	responseCode := binary.LittleEndian.Uint16(response[0x22:0x24])
	if responseCode != 0 {
		return -1, fmt.Errorf("failed to read power value code (%04x)", responseCode)
	}

	data, err := d.DecodePacket(response)
	if err != nil {
		return -1, err
	}

	// иногда приходят значения 0.FF, FFFF.FF что скорее всего является
	// ошибкой устройства поэтому игнорируем такие значения, приравнивая их к 0
	if data[0x05] == 255 {
		return 0, nil
	}

	str := fmt.Sprintf("%X.%X", uint64(data[0x07])*256+uint64(data[0x06]), data[0x05])

	return strconv.ParseFloat(str, 64)
}
