package z_stack

import (
	"context"
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
