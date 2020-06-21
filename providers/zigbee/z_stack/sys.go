package z_stack

import (
	"context"
)

const (
	ADCChannelAIN0              = 0x00
	ADCChannelAIN1              = 0x01
	ADCChannelAIN2              = 0x02
	ADCChannelAIN3              = 0x03
	ADCChannelAIN4              = 0x04
	ADCChannelAIN5              = 0x05
	ADCChannelAIN6              = 0x06
	ADCChannelAIN7              = 0x07
	ADCChannelTemperatureSensor = 0x0E
	ADCChannelVoltageReading    = 0x0F

	ADCResolutionBit8  = 0x00
	ADCResolutionBit10 = 0x01
	ADCResolutionBit12 = 0x02
	ADCResolutionBit14 = 0x03
)

type SysVersion struct {
	TransportRevision uint8
	Product           uint8
	MajorRelease      uint8
	MinorRelease      uint8
	MainTrel          uint8
	HardwareRevision  uint32
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

	data := response.DataAsBuffer()

	return &SysVersion{
		TransportRevision: data.ReadUint8(),
		Product:           data.ReadUint8(), // zStack12 = 0, zStack3x0 = 1, zStack30x = 2,
		MajorRelease:      data.ReadUint8(),
		MinorRelease:      data.ReadUint8(),
		MainTrel:          data.ReadUint8(),
		HardwareRevision:  data.ReadUint32(),
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
