package z_stack

import (
	"context"

	"github.com/kihamo/boggart/protocols/serial"
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

func (c *Client) UtilGetDeviceInfo() (*UtilDeviceInfo, error) {
	request := &Frame{}
	request.SetCommand0(0x27)

	waiter, timeout := WaiterSREQ(request)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	response, err := c.CallWithResult(ctx, request, waiter)
	if err != nil {
		return nil, err
	}

	data := response.DataAsBuffer()
	devicesCount := (data.Len() - 14) / 2

	info := &UtilDeviceInfo{
		Status:           data.ReadUint8() == 0,
		IEEEAddr:         serial.Reverse(data.Next(8)),
		ShortAddr:        data.ReadUint16(),
		DeviceType:       data.ReadUint8(),
		DeviceState:      data.ReadUint8(),
		NumAssocDevices:  data.ReadUint8(),
		AssocDevicesList: make([]uint16, 0, devicesCount),
	}

	for i := 0; i < devicesCount; i++ {
		info.AssocDevicesList = append(info.AssocDevicesList, data.ReadUint16())
	}

	return info, err
}
