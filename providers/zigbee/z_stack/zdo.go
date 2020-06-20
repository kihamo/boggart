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

/*
	ZDO_STARTUP_FROM_APP

	This command starts the device in the network.

	Usage:
		SREQ:
			       1      |       1     |       1     |     2
			Length = 0x01 | Cmd0 = 0x25 | Cmd1 = 0x40 | StartDelay
		Attributes:
			StartDelay 2 bytes Specifies the time delay before the device starts.

		SRSP:
			       1      |       1     |       1     |    1
			Length = 0x01 | Cmd0 = 0x65 | Cmd1 = 0x40 | Status
		Attributes:
			Status 1 byte 0x00 – Restored network state
			              0x01 – New network state
			              0x02 – Leave and not Started

	Example from zigbee2mqtt:
		zigbee-herdsman:adapter:zStack:znp:SREQ --> ZDO - startupFromApp - {"startdelay":100} +15ms
		zigbee-herdsman:adapter:zStack:unpi:writer --> frame [254,2,37,64,100,0,3] +14ms
		zigbee-herdsman:adapter:zStack:unpi:parser <-- [254,1,101,64,0,36] +722ms
		zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [254,1,101,64,0,36] +1ms
		zigbee-herdsman:adapter:zStack:unpi:parser --> parsed 1 - 3 - 5 - 64 - [0] - 36 +0ms
		zigbee-herdsman:adapter:zStack:znp:SRSP <-- ZDO - startupFromApp - {"status":0} +724ms
		zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [] +1ms
		zigbee-herdsman:adapter:zStack:unpi:parser <-- [254,1,69,192,9,141] +4ms
		zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [254,1,69,192,9,141] +0ms
		zigbee-herdsman:adapter:zStack:unpi:parser --> parsed 1 - 2 - 5 - 192 - [9] - 141 +0ms
*/
func (c *Client) ZDOStartupFromApp(ctx context.Context, delay uint8) (status int8, err error) {
	request := &Frame{}
	request.SetCommand0(0x25) // Type 0x1, SubSystem 0x5
	request.SetCommandID(0x40)
	request.SetData([]byte{delay})

	waiter, timeout := WaiterSREQ(request)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	response, err := c.CallWithResult(ctx, request, waiter)
	if err != nil {
		return -1, err
	}

	return int8(response.Data()[0]), nil
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

/*
	ZDO_ACTIVE_EP_REQ

	This command is generated to request a list of active endpoint from the destination device.

	Usage:
		SREQ:
			       1      |       1     |       1     |     2   |         2
			Length = 0x04 | Cmd0 = 0x25 | Cmd1 = 0x05 | DstAddr | NWKAddrOfInterest
		Attributes:
			DstAddr           2 bytes Specifies NWK address of the device generating the inquiry.
			NWKAddrOfInterest 2 bytes Specifies NWK address of the destination device being queried.

		SRSP:
			       1      |       1     |       1     |    1
			Length = 0x01 | Cmd0 = 0x65 | Cmd1 = 0x05 | Status
		Attributes:
			Status 1 byte Status is either Success (0) or Failure (1).

	Example from zigbee2mqtt:
		zigbee-herdsman:adapter:zStack:znp:SREQ --> ZDO - activeEpReq - {"dstaddr":0,"nwkaddrofinterest":0} +727ms
		zigbee-herdsman:adapter:zStack:unpi:writer --> frame [254,4,37,5,0,0,0,0,36] +728ms
		zigbee-herdsman:adapter:zStack:unpi:parser <-- [254,1,101,5,0,97] +8ms
		zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [254,1,101,5,0,97] +1ms
		zigbee-herdsman:adapter:zStack:unpi:parser --> parsed 1 - 3 - 5 - 5 - [0] - 97 +0ms
		zigbee-herdsman:adapter:zStack:znp:SRSP <-- ZDO - activeEpReq - {"status":0} +15ms
		zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [] +0ms
		zigbee-herdsman:adapter:zStack:unpi:parser <-- [254,6,69,133,0,0,0,0,0,0,198] +4ms
		zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [254,6,69,133,0,0,0,0,0,0,198] +0ms
		zigbee-herdsman:adapter:zStack:unpi:parser --> parsed 6 - 2 - 5 - 133 - [0,0,0,0,0,0] - 198 +1ms
*/
func (c *Client) ZDOActiveEndpoints(ctx context.Context) error {
	request := &Frame{}
	request.SetCommand0(0x25) // Type 0x1, SubSystem 0x5
	request.SetCommandID(0x05)
	request.SetData([]byte{
		0, 0, // DstAddr
		0, 0, // NWKAddrOfInterest
	})

	waiter, timeout := WaiterSREQ(request)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	response, err := c.CallWithResult(ctx, request, waiter)
	if err != nil {
		return err
	}

	if response.Command0() != 0x65 {
		return errors.New("bad response")
	}

	data := response.Data()
	if len(data) == 0 || data[0] != 0 {
		return errors.New("failure")
	}

	return nil
}
