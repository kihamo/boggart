package zstack

import (
	"context"
	"errors"
)

type UtilDeviceInfo struct {
	IEEEAddr         []byte
	AssocDevicesList []uint16
	ShortAddr        uint16
	Status           CommandStatus
	DeviceType       uint8
	DeviceState      DeviceState
	NumAssocDevices  uint8
}

/*
	UTIL_GET_DEVICE_INFO

	This command is sent by the tester to retrieve the device info.

	Usage:
		SREQ:
			       1      |      1      |      1
			Length = 0x00 | Cmd0 = 0x27 | Cmd1 = 0x00
		SRSP:
			       1      |      1      |      1      |   1    |    8     |     2     |      1     |      1      |        1        |     0-128
			Length = 0x02 | Cmd0 = 0x67 | Cmd1 = 0x00 | Status | IEEEAddr | ShortAddr | DeviceType | DeviceState | NumAssocDevices | AssocDeviceList
		Attributes:
			Status           1 byte      Status is a one byte field and is either success(0) or fail(1). The fail status is returned if the address value in the command message was not within the valid range.
			IEEEAddr         8 bytes     IEEE address of the device
			ShortAddr        2 bytes     Short address of the device
			DeviceType       1 byte      Indicates device type, where bits 0 to 2 indicate the capability for the device to operate as a coordinator, router, or end device, respectively.
			DeviceState      1 byte      Indicates the state of the device with different possible states as shown below:
			                             0x00: Initialized - not started automatically
			                             0x01: Initialized - not connected to anything
			                             0x02: Discovering PAN's to join
			                             0x03: Joining a PAN
			                             0x04: Rejoining a PAN, only for end devices
			                             0x05: Joined but not yet authenticated by trust center
			                             0x06: Started as device after authentication
			                             0x07: Device joined, authenticated and is a router
			                             0x08: Starting as ZigBee Coordinator
			                             0x09: Started as ZigBee Coordinator
			                             0x0A: Device has lost information about its parent
			NumAssocDevices  1 byte      Specifies the number of devices being associated to the target device.
			AssocDevicesList 0-128 bytes Array of 16-bits specifies the network address associated with the device.
*/
func (c *Client) UtilGetDeviceInfo(ctx context.Context) (*UtilDeviceInfo, error) {
	response, err := c.CallWithResultSREQ(ctx, NewFrame(0x27, CommandGetDeviceInfo))
	if err != nil {
		return nil, err
	}

	return UtilGetDeviceInfoParse(response)
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
func (c *Client) UtilLEDControl(ctx context.Context, ledID uint8, mode bool) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(ledID)  // LedId
	dataIn.WriteBoolean(mode) // Mode

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x27, CommandLEDControl))
	if err != nil {
		return err
	}

	dataOut := response.DataAsBuffer()
	if dataOut.Len() == 0 {
		return errors.New("failure")
	}

	if status := dataOut.ReadCommandStatus(); status != CommandStatusSuccess {
		return status
	}

	return nil
}
