package z_stack

import (
	"context"
)

type SysVersion struct {
	TransportRevision uint8
	Product           uint8
	MajorRelease      uint8
	MinorRelease      uint8
	MainTrel          uint8
	HardwareRevision  uint32
}

func (v *SysVersion) Type() string {
	switch v.Product {
	case VersionZStack12:
		return "zStack12"
	case VersionZStack3x0:
		return "zStack3x0"
	case VersionZStack30x:
		return "zStack30x"
	}

	return ""
}

func (c *Client) SysPing(ctx context.Context) (uint16, error) {
	request := &Frame{}
	request.SetCommand0(0x21)
	request.SetCommandID(0x01)

	response, err := c.CallWithResultSREQ(ctx, request)
	if err != nil {
		return 0, err
	}

	return response.DataAsBuffer().ReadUint16(), err
}

func (c *Client) SysVersion(ctx context.Context) (*SysVersion, error) {
	request := &Frame{}
	request.SetCommand0(0x21)
	request.SetCommandID(0x02)

	response, err := c.CallWithResultSREQ(ctx, request)
	if err != nil {
		return nil, err
	}

	dataOut := response.DataAsBuffer()

	return &SysVersion{
		TransportRevision: dataOut.ReadUint8(),
		Product:           dataOut.ReadUint8(), // zStack12 = 0, zStack3x0 = 1, zStack30x = 2,
		MajorRelease:      dataOut.ReadUint8(),
		MinorRelease:      dataOut.ReadUint8(),
		MainTrel:          dataOut.ReadUint8(),
		HardwareRevision:  dataOut.ReadUint32(),
	}, err
}

func (c *Client) SysADCRead(ctx context.Context, channel, resolution uint8) (uint16, error) {
	request := &Frame{}
	request.SetCommand0(0x21)
	request.SetCommandID(0x0D)
	request.SetData([]byte{channel, resolution})

	_, err := c.CallWithResultSREQ(ctx, request)
	if err != nil {
		return 0, err
	}

	// TODO
	//fmt.Println(response.Data())

	return 0, err
}
