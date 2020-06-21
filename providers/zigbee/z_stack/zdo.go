package z_stack

import (
	"context"
	"errors"
	"fmt"

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

type ZDOGroup struct {
	Status uint8
	ID     uint16
	Name   []byte
}

type ZDODeviceJoined struct {
	NetworkAddress uint16
	ExtendAddress  []byte
	ParentAddress  uint16
}

type ZDOEndDeviceAnnounce struct {
	SourceAddress  uint16
	NetworkAddress uint16
	IEEEAddress    []byte
	Capabilites    uint8
}

type ZDODeviceLeave struct {
	SourceAddress  uint16
	ExtendAddress  []byte
	Request        uint8
	RemoveChildren uint8
	Rejoin         uint8
}

func (c *Client) ZDOExtNwkInfo(ctx context.Context) (*ExtNwkInfo, error) {
	request := &Frame{}
	request.SetCommand0(0x25) // Type 0x1, SubSystem 0x5
	request.SetCommandID(0x50)

	response, err := c.CallWithResultSREQ(ctx, request)
	if err != nil {
		return nil, err
	}

	dataOut := response.DataAsBuffer()

	return &ExtNwkInfo{
		ShortAddr:     dataOut.ReadUint16(),
		DevState:      dataOut.ReadUint8(),
		PanID:         dataOut.ReadUint16(),
		ParentAddr:    dataOut.ReadUint16(),
		ExtendedPanID: serial.Reverse(dataOut.Next(8)),
		ParentExtAddr: serial.Reverse(dataOut.Next(8)),
		Channel:       dataOut.ReadUint8(),
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

	response, err := c.CallWithResultSREQ(ctx, request)
	if err != nil {
		return -1, err
	}

	return int8(response.Data()[0]), nil
}

/*
	Example from zigbee2mqtt:
		zigbee-herdsman:controller:log Permit joining +250ms
		zigbee-herdsman:adapter:zStack:znp:SREQ --> ZDO - mgmtPermitJoinReq - {"addrmode":15,"dstaddr":65532,"duration":254,"tcsignificance":0} +13ms
		zigbee-herdsman:adapter:zStack:unpi:writer --> frame [254,5,37,54,15,252,255,254,0,228] +13ms
		zigbee-herdsman:adapter:zStack:unpi:parser <-- [254,1,101,54,0,82] +21ms
		zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [254,1,101,54,0,82] +0ms
		zigbee-herdsman:adapter:zStack:unpi:parser --> parsed 1 - 3 - 5 - 54 - [0] - 82 +0ms
		zigbee-herdsman:adapter:zStack:znp:SRSP <-- ZDO - mgmtPermitJoinReq - {"status":0} +21ms
		zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [] +1ms
*/
func (c *Client) ZDOPermitJoin(ctx context.Context, seconds uint8) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(0x0F)    // AddrMode (поле отсутствует в оригинальной документации) networkAddress === null ? 0x0F : 0x02;
	dataIn.WriteUint16(0xFFFC) // DstAddr (0xFFFC -- broadcast)
	dataIn.WriteUint8(seconds) // Duration
	dataIn.WriteUint8(0)       // TCSignificance

	request := &Frame{}
	request.SetCommand0(0x25) // Type 0x1, SubSystem 0x5
	request.SetCommandID(0x36)
	request.SetDataAsBuffer(dataIn)

	response, err := c.CallWithResultSREQ(ctx, request)
	if err != nil {
		return err
	}

	if response.Command0() != 0x65 && response.Command1() != 0x36 {
		return errors.New("bad response")
	}

	dataOut := response.Data()
	if len(dataOut) == 0 || dataOut[0] != 0 {
		return errors.New("failure")
	}

	return nil
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
	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(0) // DstAddr
	dataIn.WriteUint16(0) // NWKAddrOfInterest

	request := &Frame{}
	request.SetCommand0(0x25) // Type 0x1, SubSystem 0x5
	request.SetCommandID(0x05)
	request.SetDataAsBuffer(dataIn)

	response, err := c.CallWithResultSREQ(ctx, request)
	if err != nil {
		return err
	}

	if response.Command0() != 0x65 {
		return errors.New("bad response")
	}

	dataOut := response.Data()
	if len(dataOut) == 0 || dataOut[0] != 0 {
		return errors.New("failure")
	}

	return nil
}

func (c *Client) ZDONodeDescription(ctx context.Context, DstAddr, NWKAddrOfInterest uint16) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(DstAddr)           // DstAddr
	dataIn.WriteUint16(NWKAddrOfInterest) // NWKAddrOfInterest

	request := &Frame{}
	request.SetCommand0(0x25) // Type 0x1, SubSystem 0x5
	request.SetCommandID(0x02)
	request.SetDataAsBuffer(dataIn)

	response, err := c.CallWithResultSREQ(ctx, request)
	if err != nil {
		return err
	}

	fmt.Println(response)

	return nil
}

/*
	Usage:
		SREQ:
			       1      |       1     |       1     |     1    |    2
			Length = 0x03 | Cmd0 = 0x25 | Cmd1 = 0x4A | Endpoint | GroupID
		Attributes:
			Endpoint 1 byte  Endpoint ID.
			GroupID  2 bytes Group ID.

		SRSP:
			       1      |       1     |       1     |    1   |    2    |    1    |     ?
			Length = 0x01 | Cmd0 = 0x25 | Cmd1 = 0x02 | Status | GroupID | NameLen | GroupName
		Attributes:
			Status    1 byte  Status is either exist (0) or not exists (1).
			GroupID   2 bytes Group ID.
			NameLen   1 byte  Group name length.
			GroupName ? bytes Group name

	Example from zigbee2mqtt:
		zigbee-herdsman:adapter:zStack:znp:SREQ --> ZDO - extFindGroup - {"endpoint":242,"groupid":2948} +7ms
		zigbee-herdsman:adapter:zStack:unpi:writer --> frame [254,3,37,74,242,132,11,17] +6ms
		zigbee-herdsman:adapter:zStack:unpi:parser <-- [254,19,101,74,0,132,11,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,179] +11ms
		zigbee-herdsman:adapter:zStack:unpi:parser --- parseNext [254,19,101,74,0,132,11,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,179] +0ms
		zigbee-herdsman:adapter:zStack:unpi:parser --> parsed 19 - 3 - 5 - 74 - [0,132,11,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0] - 179 +1ms
		zigbee-herdsman:adapter:zStack:znp:SRSP <-- ZDO - extFindGroup - {"status":0,"groupid":2948,"namelen":0,"groupname":{"type":"Buffer","data":[]}} +13ms
*/
func (c *Client) ZDOExtFindGroup(ctx context.Context, endpoint uint8, group uint16) (*ZDOGroup, error) {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(endpoint)
	dataIn.WriteUint16(group)

	request := &Frame{}
	request.SetCommand0(0x25) // Type 0x1, SubSystem 0x5
	request.SetCommandID(0x4A)
	request.SetDataAsBuffer(dataIn)

	response, err := c.CallWithResultSREQ(ctx, request)
	if err != nil {
		return nil, err
	}

	dataOut := response.DataAsBuffer()
	if dataOut.Len() < 4 {
		return nil, errors.New("failure")
	}

	g := &ZDOGroup{
		Status: dataOut.ReadUint8(),
		ID:     dataOut.ReadUint16(),
	}
	dataOut.ReadUint8()      // namelen
	g.Name = dataOut.Bytes() // name

	return g, nil
}

func (c *Client) ZDOExtAddToGroup(ctx context.Context, endpoint uint8, group uint16, groupName []byte) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(endpoint)              // endpoint
	dataIn.WriteUint16(group)                // groupid
	dataIn.WriteUint8(uint8(len(groupName))) // namelen
	dataIn.Write(groupName)                  // groupname

	request := &Frame{}
	request.SetCommand0(0x25) // Type 0x1, SubSystem 0x5
	request.SetCommandID(0x02)
	request.SetDataAsBuffer(dataIn)

	response, err := c.CallWithResultSREQ(ctx, request)
	if err != nil {
		return err
	}

	fmt.Println("ZDOExtAddToGroup", response)

	return nil
}

func (c *Client) ZDODeviceJoinedMessage(frame *Frame) (*ZDODeviceJoined, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != 0xCA {
		return nil, errors.New("frame isn't a device joined command")
	}

	dataOut := frame.DataAsBuffer()

	return &ZDODeviceJoined{
		NetworkAddress: dataOut.ReadUint16(),
		ExtendAddress:  dataOut.ReadIEEEAddr(),
		ParentAddress:  dataOut.ReadUint16(),
	}, nil
}

/*
	ZDO_END_DEVICE_ANNCE_IND

	This callback indicates the ZDO End Device Announce.

	Usage:
		AREQ:
			       1      |      1      |      1      |    2    |    2    |    8     |      1
			Length = 0x0D | Cmd0 = 0x45 | Cmd1 = 0xC1 | SrcAddr | NwkAddr | IEEEAddr | Capabilites
		Attributes:
			SrcAddr     2 bytes Source address of the message.
			NwkAddr     2 bytes Specifies the device’ s short address.
			IEEEAddr    8 bytes Specifies the 64 bit IEEE address of source device.
			Capabilites 1 byte  Specifies the MAC capabilities of the device.
			                    Bit: 0 – Alternate PAN Coordinator
			                         1 – Device type: 1- ZigBee Router; 0 – End Device
			                         2 – Power Source: 1 Main powered
			                         3 – Receiver on when Idle
			                         4 – Reserved
			                         5 – Reserved
			                         6 – Security capability
			                         7 – Reserved
*/
func (c *Client) ZDOEndDeviceAnnounceMessage(frame *Frame) (*ZDOEndDeviceAnnounce, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != 0xC1 {
		return nil, errors.New("frame isn't a end device announce command")
	}

	dataOut := frame.DataAsBuffer()

	return &ZDOEndDeviceAnnounce{
		SourceAddress:  dataOut.ReadUint16(),
		NetworkAddress: dataOut.ReadUint16(),
		IEEEAddress:    dataOut.ReadIEEEAddr(),
		Capabilites:    dataOut.ReadUint8(),
	}, nil
}

func (c *Client) ZDODeviceLeaveMessage(frame *Frame) (*ZDODeviceLeave, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != 0xC9 {
		return nil, errors.New("frame isn't a device leave command")
	}

	dataOut := frame.DataAsBuffer()

	return &ZDODeviceLeave{
		SourceAddress:  dataOut.ReadUint16(),
		ExtendAddress:  dataOut.ReadIEEEAddr(),
		Request:        dataOut.ReadUint8(),
		RemoveChildren: dataOut.ReadUint8(),
		Rejoin:         dataOut.ReadUint8(),
	}, nil
}
