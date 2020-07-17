package z_stack

import (
	"context"
	"errors"
)

type ZbConfigurationMessage struct {
	Status   CommandStatus
	ConfigID uint8
	Len      uint8
	Value    []byte
}

/*
ZB_READ_CONFIGURATION

This command is used to get a configuration property from nonvolatile memory.

Usage:
	SREQ:
		       1      |      1      |      1      |    1
		Length = 0x01 | Cmd0 = 0x26 | Cmd1 = 0x04 | ConfigId
	Attributes:
		ConfigId 1 byte Specifies the Identifier for the configuration property.

	SRSP:
		         1         |      1      |      1      |    1   |     1    |  1  | 0-128
		Length = 0x03-0x83 | Cmd0 = 0x66 | Cmd1 = 0x04 | Status | ConfigId | Len | Value
	Attributes:
		Status   1 byte      This field indicates either SUCCESS (0) or FAILURE (1).
		ConfigId 1 byte      Specifies the Identifier for the configuration property.
		Len      1 byte      Specifies the size of the Value buffer in bytes.
		Value    0-128 bytes buffer to hold the configuration property.
*/
func (c *Client) ZbReadConfiguration(ctx context.Context, configID uint8) (*ZbConfigurationMessage, error) {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(configID) // config id

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x26, 0x04))
	if err != nil {
		return nil, err
	}

	dataOut := response.DataAsBuffer()
	if dataOut.Len() == 0 {
		return nil, errors.New("failure")
	}

	msg := &ZbConfigurationMessage{
		Status:   dataOut.ReadCommandStatus(),
		ConfigID: dataOut.ReadUint8(),
		Len:      dataOut.ReadUint8(),
	}
	msg.Value = dataOut.Next(int(msg.Len))

	if msg.Status != CommandStatusSuccess {
		return nil, msg.Status
	}

	return msg, nil
}

/*
ZB_WRITE_CONFIGURATION

This command is used to write a Configuration Property to nonvolatile memory.

Usage:
	SREQ:
		         1         |      1      |      1      |    1     |  1  | 1-128
		Length = 0x03-0x83 | Cmd0 = 0x26 | Cmd1 = 0x05 | ConfigId | Len | Value
	Attributes:
		ConfigId 1 byte      The Identifier for the configuration property
		Len      1 byte      Specifies the size of the Value buffer in bytes.
		Value    1-128 bytes The buffer containing the new value of the configuration property.

	SRSP:
		       1      |      1      |      1      |    1
		Length = 0x01 | Cmd0 = 0x66 | Cmd1 = 0x05 | Status
	Attributes:
		Status 1 byte This field indicates either SUCCESS (0) or FAILURE (1).
*/
func (c *Client) ZbWriteConfiguration(ctx context.Context, configID uint8, value []byte) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(configID)          // config id
	dataIn.WriteUint8(uint8(len(value))) // len
	dataIn.Write(value)                  // value

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x26, 0x05))
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
