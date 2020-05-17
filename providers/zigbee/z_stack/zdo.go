package z_stack

import (
	"context"

	"github.com/kihamo/boggart/protocols/serial"
)

type ExtNwkInfo struct {
	ShortAddr     uint16
	DevState      uint8
	PanID         uint16
	ParentAddr    uint16
	ExtendedPanID []byte
	ParentExtAddr []byte
	Channel       uint8
}

func (c *Client) ZDOExtNwkInfo() (*ExtNwkInfo, error) {
	request := &Frame{}
	request.SetCommand0(0x25) // Type 0x1, SubSystem 0x5
	request.SetCommandID(0x50)

	waiter, timeout := WaiterSREQ(request)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	response, err := c.CallWithResult(ctx, request, waiter)
	if err != nil {
		return nil, err
	}

	data := response.DataAsBuffer()

	return &ExtNwkInfo{
		ShortAddr:     data.ReadUint16(),
		DevState:      data.ReadUint8(),
		PanID:         data.ReadUint16(),
		ParentAddr:    data.ReadUint16(),
		ExtendedPanID: serial.Reverse(data.Next(8)),
		ParentExtAddr: serial.Reverse(data.Next(8)),
		Channel:       data.ReadUint8(),
	}, nil
}
