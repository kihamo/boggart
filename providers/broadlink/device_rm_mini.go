package broadlink

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	"github.com/kihamo/boggart/providers/broadlink/internal"
)

var (
	ErrSignalNotCaptured = errors.New("signal not captured")
)

type RMMini struct {
	*internal.Device
}

func NewRMMini(mac net.HardwareAddr, addr string) *RMMini {
	return &RMMini{
		Device: internal.NewDevice(KindRMMini, mac, addr),
	}
}

func (d *RMMini) StartCaptureRemoteControlCode() error {
	/*
		Offset    Contents
		0x00      0x03
		0x01-0x0f 0x00
	*/

	payload := make([]byte, 0x10)
	payload[0] = 0x03

	response, err := d.Call(0x6a, payload)
	if err != nil {
		return err
	}

	responseCode := binary.LittleEndian.Uint16(response[0x22:0x24])
	if responseCode != 0 {
		return errors.New("failed to start capturing remote control code")
	}

	return nil
}

func (d *RMMini) ReadCapturedRemoteControlCodeRaw() ([]byte, error) {
	/*
		Offset    Contents
		0x00      0x04
		0x01-0x0f 0x00
	*/

	payload := make([]byte, 0x10)
	payload[0] = 4

	response, err := d.Call(0x6a, payload)
	if err != nil {
		return nil, err
	}

	responseCode := binary.LittleEndian.Uint16(response[0x22:0x24])
	if responseCode != 0 {
		if responseCode == 0xfff6 {
			return nil, ErrSignalNotCaptured
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

	return data[4:], nil

	//sz := int(RemoteType(binary.LittleEndian.Uint16(data[6:8])))
	//if len(data) < 8+sz {
	//	return nil, errors.New("incomplete data")
	//}

	//return data[8:8+sz], nil
}

func (d *RMMini) ReadCapturedRemoteControlCodeRawAsString() (string, error) {
	code, err := d.ReadCapturedRemoteControlCodeRaw()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%04x", code), nil
}

func (d *RMMini) ReadCapturedRemoteControlCode() (RemoteType, []byte, error) {
	code, err := d.ReadCapturedRemoteControlCodeRaw()
	if err != nil {
		return -1, nil, err
	}

	if len(code) < 4 {
		return -1, nil, errors.New("wrong length of code")
	}

	return RemoteType(binary.LittleEndian.Uint16(code[0:2])), code[4:], nil
}

func (d *RMMini) ReadCapturedRemoteControlCodeAsString() (RemoteType, string, error) {
	remoteType, code, err := d.ReadCapturedRemoteControlCode()
	if err != nil {
		return remoteType, "", err
	}

	return remoteType, fmt.Sprintf("%04x", code), err
}

func (d *RMMini) SendRemoteControlCode(remoteType RemoteType, code []byte, count int) error {
	/*
		Offset    Contents
		0x00      0x02
		0x01-0x03 0x00
		0x04      0x26 = IR, 0xb2 for RF 433Mhz, 0xd7 for RF 315Mhz
		0x05      repeat count, (0 = no repeat, 1 send twice, .....)
		0x06-0x07 Length of the following data in little endian
		0x08 .... Pulse lengths in 2^-15 s units (µs * 269 / 8192 works very well)
		....      0x0d 0x05 at the end for IR only
	*/
	if count < 0 {
		return errors.New("count must be a positive integer")
	}

	payload := make([]byte, 0x08+len(code))
	payload[0] = 0x02
	payload[4] = byte(remoteType)
	payload[5] = byte(count)

	binary.LittleEndian.PutUint16(payload[6:], uint16(len(code)))
	copy(payload[8:], code)

	response, err := d.Call(0x6a, payload)
	if err != nil {
		return err
	}

	responseCode := binary.LittleEndian.Uint16(response[0x22:0x24])
	if responseCode != 0 {
		return fmt.Errorf("failed sending remote control code (%04x)", responseCode)
	}

	return nil
}

func (d *RMMini) SendRemoteControlCodeRaw(code []byte, count int) error {
	if len(code) < 2 {
		return errors.New("wrong length of code")
	}

	remoteType := RemoteType(binary.LittleEndian.Uint16(code[0:2]))
	switch remoteType {
	case RemoteIR, RemoteRF433Mhz, RemoteRF315Mhz:
	default:
		return fmt.Errorf("unknown remote type (%04x)", remoteType)
	}

	return d.SendRemoteControlCode(remoteType, code[4:], count)
}

func (d *RMMini) SendRemoteControlCodeRawAsString(code string, count int) error {
	decoded, err := hex.DecodeString(code)
	if err != nil {
		return err
	}

	return d.SendRemoteControlCodeRaw(decoded, count)
}

func (d *RMMini) SendIRRemoteControlCode(code []byte, count int) error {
	return d.SendRemoteControlCode(RemoteIR, code, count)
}

func (d *RMMini) SendIRRemoteControlCodeAsString(code string, count int) error {
	decoded, err := hex.DecodeString(code)
	if err != nil {
		return err
	}

	return d.SendRemoteControlCode(RemoteIR, decoded, count)
}
