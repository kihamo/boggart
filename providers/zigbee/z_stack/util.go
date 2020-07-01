package z_stack

import (
	"context"
	"errors"
)

type UtilDeviceInfo struct {
	Status           bool
	IEEEAddr         []byte
	ShortAddr        uint16
	DeviceType       uint8
	DeviceState      uint8
	NumAssocDevices  uint8
	AssocDevicesList []uint16
}

func (c *Client) UtilGetDeviceInfo(ctx context.Context) (*UtilDeviceInfo, error) {
	response, err := c.CallWithResultSREQ(ctx, NewFrame(0x27, 0x00))
	if err != nil {
		return nil, err
	}

	dataOut := response.DataAsBuffer()
	devicesCount := (dataOut.Len() - 14) / 2

	info := &UtilDeviceInfo{
		Status:           dataOut.ReadUint8() == 0,
		IEEEAddr:         dataOut.ReadIEEEAddr(),
		ShortAddr:        dataOut.ReadUint16(),
		DeviceType:       dataOut.ReadUint8(),
		DeviceState:      dataOut.ReadUint8(),
		NumAssocDevices:  dataOut.ReadUint8(),
		AssocDevicesList: make([]uint16, 0, devicesCount),
	}

	for i := 0; i < devicesCount; i++ {
		info.AssocDevicesList = append(info.AssocDevicesList, dataOut.ReadUint16())
	}

	return info, err
}

/**
UTIL_LED_CONTROL

This command is used by the tester to control the LEDs on the board.

Usage:
	SREQ:
		       1      |      1      |      1      |   1   |  1
		Length = 0x02 | Cmd0 = 0x27 | Cmd1 = 0x0A | LedId | Mode
	Attributes:
		Laded 1 byte The LED number
		Mode  1 byte 0: OFF, 1: ON

	SRSP:
		       1      |      1      |      1      |    1
		Length = 0x01 | Cmd0 = 0x67 | Cmd1 = 0x0A | Status
	Attributes:
		Status 1 byte Status is either Success (0) or Failure (1).
*/
func (c *Client) UtilLEDControl(ctx context.Context, LedId uint8, mode bool) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(LedId)  // LedId
	dataIn.WriteBoolean(mode) // Mode

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x27, 0x0A))
	if err != nil {
		return err
	}

	dataOut := response.DataAsBuffer()
	if dataOut.Len() == 0 || dataOut.ReadUint8() != 0 {
		return errors.New("failure")
	}

	return nil
}
