package zstack

import (
	"context"
	"errors"
)

type SysVersion struct {
	TransportRevision uint8
	Product           uint8
	MajorRelease      uint8
	MinorRelease      uint8
	MainTrel          uint8
	HardwareRevision  uint32
}

type NVItem struct {
	Status CommandStatus
	Value  []byte
}

func (i *NVItem) DataAsBuffer() *Buffer {
	return NewBuffer(i.Value)
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
	response, err := c.CallWithResultSREQ(ctx, NewFrame(0x21, 0x01))
	if err != nil {
		return 0, err
	}

	return response.DataAsBuffer().ReadUint16(), err
}

func (c *Client) SysVersion(ctx context.Context) (*SysVersion, error) {
	response, err := c.CallWithResultSREQ(ctx, NewFrame(0x21, 0x02))
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
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(channel)    // channel
	dataIn.WriteUint8(resolution) // resolution

	_, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x21, 0x0D))
	if err != nil {
		return 0, err
	}

	// TODO
	//fmt.Println(response.Data())

	return 0, err
}

func (c *Client) SysOsalNvItemInit(ctx context.Context, id uint16, value []byte) error {
	l := len(value)

	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(id)
	dataIn.WriteUint16(uint16(l))
	dataIn.WriteUint8(uint8(l))
	dataIn.Write(value)

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x21, 0x07))
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

/**
SYS_OSAL_NV_READ

This command is used by the tester to read a single memory item in the target non-volatile memory. The command accepts an attribute Id value and returns the memory value present in the target for the specified attribute Id.

Usage:
	SREQ:
		       1      |      1      |      1      |  2 |   1
		Length = 0x03 | Cmd0 = 0x21 | Cmd1 = 0x08 | Id | Offset
	Attributes:
		Id     2 bytes The Id of the NV item.
		Offset 1 byte  Number of bytes offset from the beginning or the NV value.

	SRSP:
		         1         |      1      |      1      |    1   |  1  | 0-128
		Length = 0x02-0x82 | Cmd0 = 0x61 | Cmd1 = 0x08 | Status | Len | Value
	Attributes:
		Status 1 byte      Status is either Success (0) or Failure (1).
		Len    1 byte      Length of the NV value.
		Value  0-128 bytes Value of the NV item.
*/
func (c *Client) SysOsalNvRead(ctx context.Context, id uint16, offset uint8) (*NVItem, error) {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(id)    // id
	dataIn.WriteUint8(offset) // offset

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x21, 0x08))
	if err != nil {
		return nil, err
	}

	dataOut := response.DataAsBuffer()
	item := &NVItem{
		Status: dataOut.ReadCommandStatus(),
	}
	l := dataOut.ReadUint8()
	item.Value = dataOut.Next(int(l))

	if item.Status != CommandStatusSuccess && item.Status != CommandStatusInvalidParam {
		return nil, item.Status
	}

	return item, nil
}

/**
SYS_OSAL_NV_WRITE

This command is used by the tester to write to a particular item in non-volatile memory. The command accepts an attribute Id and an attribute value. The attribute value is written to the location specified for the attribute Id in the target.

Usage:
	SREQ:
		         1         |      1      |      1      |  2 |   1    |  1  | 1-128
		Length = 0x04-0x84 | Cmd0 = 0x21 | Cmd1 = 0x09 | Id | Offset | Len | Value
	Attributes:
		Id     2 bytes     The Id of the NV item.
		Offset 1 byte      Number of bytes offset from the beginning or the NV value.
		Len    1 byte      Length of the NV value.
		Value  0-128 bytes Value of the NV item.

	SRSP:
		       1      |      1      |      1      |    1
		Length = 0x01 | Cmd0 = 0x61 | Cmd1 = 0x09 | Status
	Attributes:
		Status 1 byte Status is either Success (0) or Failure (1).

*/
func (c *Client) SysOsalNvWrite(ctx context.Context, id uint16, offset uint8, value []byte) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(id)               // id
	dataIn.WriteUint8(offset)            // offset
	dataIn.WriteUint8(uint8(len(value))) // length
	dataIn.Write(value)                  // value

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x21, 0x09))
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

func (c *Client) SysReset(ctx context.Context, t uint8) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(t) // type

	return c.Call(dataIn.Frame(0x41, 0x00))
}
