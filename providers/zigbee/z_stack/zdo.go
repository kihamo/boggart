package z_stack

import (
	"context"
	"errors"

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

func (c *Client) ZDOExtNwkInfo(ctx context.Context) (*ExtNwkInfo, error) {
	request := &Frame{}
	request.SetCommand0(0x25) // Type 0x1, SubSystem 0x5
	request.SetCommandID(0x50)

	waiter, timeout := WaiterSREQ(request)
	ctx, cancel := context.WithTimeout(ctx, timeout)
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

func (c *Client) ZDOPermitJoin(ctx context.Context, seconds uint8) (interface{}, error) {
	request := &Frame{}
	request.SetCommand0(0x25) // Type 0x1, SubSystem 0x5
	request.SetCommandID(0x36)
	request.SetData([]byte{
		0xFF, 0xFC, // broadcast
		seconds, 0,
	})

	waiter, timeout := WaiterSREQ(request)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	response, err := c.CallWithResult(ctx, request, waiter)
	if err != nil {
		return nil, err
	}

	if response.Command0() != 0x65 && response.Command1() != 0x36 {
		return nil, errors.New("bad response")
	}

	data := response.Data()
	if len(data) == 0 || data[0] != 0 {
		return nil, errors.New("failure")
	}

	return nil, nil
}
