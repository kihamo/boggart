package z_stack

import (
	"context"
	"encoding/binary"

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
	request.SetCommandID(0x00)

	waiter, timeout := WaiterSREQ(request)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	response, err := c.CallWithResult(ctx, request, waiter)
	if err != nil {
		return nil, err
	}

	data := response.Data()
	devicesCount := (len(data) - 14) / 2

	info := &UtilDeviceInfo{
		Status:           data[0] == byte(0),
		IEEEAddr:         serial.Reverse(data[1:9]),
		ShortAddr:        binary.LittleEndian.Uint16(data[9:11]),
		DeviceType:       data[11],
		DeviceState:      data[12],
		NumAssocDevices:  data[13],
		AssocDevicesList: make([]uint16, 0, devicesCount),
	}

	for i := 0; i < devicesCount; i++ {
		info.AssocDevicesList = append(info.AssocDevicesList,
			binary.LittleEndian.Uint16(data[14+(i*2):14+2+(i*2)]),
		)
	}

	return info, err
}
