package z_stack

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/protocols/serial"
)

type ZDOExtNetworkInfo struct {
	ShortAddr     uint16
	DevState      uint8
	PanID         uint16
	ParentAddr    uint16
	ExtendedPanID []byte
	ParentExtAddr []byte
	Channel       uint8
}

type ZDOGroup struct {
	Status CommandStatus
	ID     uint16
	Name   []byte
}

func (c *Client) ZDOExtNetworkInfo(ctx context.Context) (*ZDOExtNetworkInfo, error) {
	response, err := c.CallWithResultSREQ(ctx, NewFrame(0x25, 0x50))
	if err != nil {
		return nil, err
	}

	dataOut := response.DataAsBuffer()

	return &ZDOExtNetworkInfo{
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
func (c *Client) ZDOStartupFromApp(ctx context.Context, delay uint8) (uint8, error) {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(delay)

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x40))
	if err != nil {
		return 0, err
	}

	return response.DataAsBuffer().ReadUint8(), nil
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

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x36))
	if err != nil {
		return err
	}

	if response.Command0() != 0x65 && response.Command1() != 0x36 {
		return errors.New("bad response")
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

/*
	ZDO_SIMPLE_DESC_REQ

	This command is generated to inquire as to the Simple Descriptor of the destination device’s Endpoint.

	Usage:
		SREQ:
			       1      |      1      |      1      |    2    |         2         |     1
			Length = 0x05 | Cmd0 = 0x25 | Cmd1 = 0x04 | DstAddr | NWKAddrOfInterest | Endpoint
		Attributes:
			DstAddr           2 byte  Specifies NWK address of the device generating the inquiry.
			NWKAddrOfInterest 2 bytes Specifies NWK address of the destination device being queried.
			Endpoint          1 byte  Specifies the application endpoint the data is from.
		SRSP:
			       1      |      1      |      1      |    1
			Length = 0x01 | Cmd0 = 0x65 | Cmd1 = 0x04 | Status
		Attributes:
			Status 1 byte Status is either Success (0) or Failure (1).
*/
func (c *Client) ZDOSimpleDescriptor(ctx context.Context, dstAddress, networkAddress uint16, endpoint uint8) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(dstAddress)     // DstAddr
	dataIn.WriteUint16(networkAddress) // NWKAddrOfInterest
	dataIn.WriteUint8(endpoint)        // Endpoint

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x04))
	if err != nil {
		return err
	}

	if response.Command0() != 0x65 {
		return errors.New("bad response")
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
func (c *Client) ZDOActiveEndpoints(ctx context.Context, dstAddress, networkAddress uint16) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(dstAddress)     // DstAddr
	dataIn.WriteUint16(networkAddress) // NWKAddrOfInterest

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x05))
	if err != nil {
		return err
	}

	if response.Command0() != 0x65 {
		return errors.New("bad response")
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

func (c *Client) ZDONodeDescription(ctx context.Context, DstAddr, NWKAddrOfInterest uint16) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(DstAddr)           // DstAddr
	dataIn.WriteUint16(NWKAddrOfInterest) // NWKAddrOfInterest

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x02))
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

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x4A))
	if err != nil {
		return nil, err
	}

	dataOut := response.DataAsBuffer()
	if dataOut.Len() < 4 {
		return nil, errors.New("failure")
	}

	g := &ZDOGroup{
		Status: dataOut.ReadCommandStatus(),
		ID:     dataOut.ReadUint16(),
	}
	l := dataOut.ReadUint8()      // namelen
	g.Name = dataOut.Next(int(l)) // name

	return g, nil
}

func (c *Client) ZDOExtAddToGroup(ctx context.Context, endpoint uint8, group uint16, groupName []byte) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(endpoint)              // endpoint
	dataIn.WriteUint16(group)                // groupid
	dataIn.WriteUint8(uint8(len(groupName))) // namelen
	dataIn.Write(groupName)                  // groupname

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x02))
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

func (c *Client) ZDOLQI(ctx context.Context, networkAddress uint16, startIndex uint8) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(networkAddress) // DstAddr
	dataIn.WriteUint8(startIndex)      // StartIndex

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x31))
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

/*
	ZDO_MGMT_NWK_DISC_REQ

	This command is generated to request the destination device to perform a network discovery.

	Usage:
		SREQ:
			       1       |      1      |       1     |    2    |       4      |      1       |      1
			Length =  0x08 | Cmd0 = 0x45 | Cmd1 = 0x30 | DstAddr | ScanChannels | ScanDuration | StartIndex
		Attributes:
			DstAddr      2 bytes Specifies the network address of the device performing the discovery.
			ScanChannels 4 bytes Specifies the Bit Mask for channels to scan:
			                     NONE         0x00000000
			                     ALL_CHANNELS 0x07FFF800
			                     CHANNEL 11   0x00000800
			                     CHANNEL 12   0x00001000
			                     CHANNEL 13   0x00002000
			                     CHANNEL 14   0x00004000
			                     CHANNEL 15   0x00008000
			                     CHANNEL 16   0x00010000
			                     CHANNEL 17   0x00020000
			                     CHANNEL 18   0x00040000
			                     CHANNEL 19   0x00080000
			                     CHANNEL 20   0x00100000
			                     CHANNEL 21   0x00200000
			                     CHANNEL 22   0x00400000
			                     CHANNEL 23   0x00800000
			                     CHANNEL 24   0x01000000
			                     CHANNEL 25   0x02000000
			                     CHANNEL 26   0x04000000
			ScanDuration 1 byte  Specifies the scanning time.
			StartIndex   1 byte  Specifies where to start in the response array list. The result may contain more entries than can be reported, so this field allows the user to retrieve the responses anywhere in the array list.
		SRSP:
			       1      |      1      |      1      |   1
			Length = 0x01 | Cmd0 = 0x65 | Cmd1 = 0x30 | Status
		Attributes:
			Status 1 byte Status is either Success (0) or Failure (1).
*/
func (c *Client) ZDOManagementNetworkDiscovery(ctx context.Context, dstAddr uint16, scanChannels uint32, scanDuration, startIndex uint8) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(dstAddr)      // DstAddr
	dataIn.WriteUint32(scanChannels) // ScanChannels
	dataIn.WriteUint8(scanDuration)  // ScanDuration
	dataIn.WriteUint8(startIndex)    // StartIndex

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x65, 0x30))
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

/*
	ZDO_MGMT_RTG_REQ

	This command is generated to request the Routing Table of the destination device

	Usage:
		SREQ:
			       1      |      1      |      1      |    2    |     1
			Length = 0x03 | Cmd0 = 0x25 | Cmd1 = 0x32 | DstAddr | StartIndex
		Attributes:
			DstAddr    2 bytes Specifies the network address of the device generating the query.
			StartIndex 1 byte  Specifies where to start in the response array list. The result may contain more entries than can be reported, so this field allows the user to retrieve the responses anywhere in the array list.

		SRSP:
			       1      |      1      |      1      |    1
			Length = 0x01 | Cmd0 = 0x65 | Cmd1 = 0x32 | Status
		Attributes:
			Status 1 byte Status is either Success (0) or Failure (1).
*/
func (c *Client) ZDORoutingTable(ctx context.Context, dstAddr uint16, startIndex uint8) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(dstAddr)   // dstAddr
	dataIn.WriteUint8(startIndex) // StartIndex

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x32))
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

func (c *Client) ZDODiscoverRoute(ctx context.Context, dstAddr uint16, options, radius uint8) error {
	dataIn := NewBuffer(nil)
	dataIn.WriteUint16(dstAddr) // dstAddr
	dataIn.WriteUint8(options)  // options
	dataIn.WriteUint8(radius)   // radius

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x45))
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
