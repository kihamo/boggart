package z_stack

import (
	"context"
	"errors"
	"fmt"

	"github.com/kihamo/boggart/protocols/serial"
)

type ExtNetworkInfo struct {
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

type ZDODeviceJoinedMessage struct {
	NetworkAddress uint16
	ExtendAddress  []byte
	ParentAddress  uint16
}

type ZDOEndDeviceAnnounceMessage struct {
	SourceAddress  uint16
	NetworkAddress uint16
	IEEEAddress    []byte
	Capabilities   uint8
}

type ZDODeviceLeaveMessage struct {
	SourceAddress  uint16
	ExtendAddress  []byte
	Request        uint8
	RemoveChildren uint8
	Rejoin         uint8
}

type NeighborLqiListItem struct {
	ExtendedPanID   []byte
	ExtendedAddress []byte
	NetworkAddress  uint16
	DeviceType      uint8
	RxOnWhenIdle    uint8
	Relationship    uint8
	PermitJoining   uint8
	Depth           uint8
	LQI             uint8
}

type ZDOLQIMessage struct {
	SourceAddress        uint16
	Status               uint8
	NeighborTableEntries uint8
	StartIndex           uint8
	NeighborLQIListCount uint8
	NeighborLqiList      []NeighborLqiListItem
}

type NetworkListItem struct {
	PAN             uint16
	LogicalChannel  uint8
	StackProfile    uint8
	ZigBeeVersion   uint8
	BeaconOrder     uint8
	SuperFrameOrder uint8
	PermitJoining   bool
}

type ZDOManagementNetworkDiscoveryMessage struct {
	SourceAddress    uint16
	Status           uint8
	NetworkCount     uint8
	StartIndex       uint8
	NetworkListCount uint8
	NetworkList      []NetworkListItem
}

func (c *Client) ZDOExtNetworkInfo(ctx context.Context) (*ExtNetworkInfo, error) {
	response, err := c.CallWithResultSREQ(ctx, NewFrame(0x25, 0x50))
	if err != nil {
		return nil, err
	}

	dataOut := response.DataAsBuffer()

	return &ExtNetworkInfo{
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
	dataIn := NewBuffer(nil)
	dataIn.WriteUint8(delay)

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x40))
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

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x36))
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

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x05))
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

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x02))
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

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x4A))
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

	response, err := c.CallWithResultSREQ(ctx, dataIn.Frame(0x25, 0x02))
	if err != nil {
		return err
	}

	fmt.Println("ZDOExtAddToGroup", response)

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

	dataOut := response.Data()
	if len(dataOut) == 0 || dataOut[0] != 0 {
		return errors.New("failure")
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

	dataOut := response.Data()
	if len(dataOut) == 0 || dataOut[0] != 0 {
		return errors.New("failure")
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

	dataOut := response.Data()
	if len(dataOut) == 0 || dataOut[0] != 0 {
		return errors.New("failure")
	}

	return nil
}

func (c *Client) ZDODeviceJoinedMessage(frame *Frame) (*ZDODeviceJoinedMessage, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != CommandTcDeviceInd {
		return nil, errors.New("frame isn't a device joined command")
	}

	dataOut := frame.DataAsBuffer()

	return &ZDODeviceJoinedMessage{
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
			Length = 0x0D | Cmd0 = 0x45 | Cmd1 = 0xC1 | SrcAddr | NwkAddr | IEEEAddr | Capabilities
		Attributes:
			SrcAddr     2 bytes Source address of the message.
			NwkAddr     2 bytes Specifies the device’ s short address.
			IEEEAddr    8 bytes Specifies the 64 bit IEEE address of source device.
			Capabilities 1 byte  Specifies the MAC capabilities of the device.
			                    Bit: 0 – Alternate PAN Coordinator
			                         1 – Device type: 1- ZigBee Router; 0 – End Device
			                         2 – Power Source: 1 Main powered
			                         3 – Receiver on when Idle
			                         4 – Reserved
			                         5 – Reserved
			                         6 – Security capability
			                         7 – Reserved
*/
func (c *Client) ZDOEndDeviceAnnounceMessage(frame *Frame) (*ZDOEndDeviceAnnounceMessage, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != CommandEndDeviceAnnounceInd {
		return nil, errors.New("frame isn't a end device announce command")
	}

	dataOut := frame.DataAsBuffer()

	return &ZDOEndDeviceAnnounceMessage{
		SourceAddress:  dataOut.ReadUint16(),
		NetworkAddress: dataOut.ReadUint16(),
		IEEEAddress:    dataOut.ReadIEEEAddr(),
		Capabilities:   dataOut.ReadUint8(),
	}, nil
}

func (c *Client) ZDODeviceLeaveMessage(frame *Frame) (*ZDODeviceLeaveMessage, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != CommandLeaveInd {
		return nil, errors.New("frame isn't a device leave command")
	}

	dataOut := frame.DataAsBuffer()

	return &ZDODeviceLeaveMessage{
		SourceAddress:  dataOut.ReadUint16(),
		ExtendAddress:  dataOut.ReadIEEEAddr(),
		Request:        dataOut.ReadUint8(),
		RemoveChildren: dataOut.ReadUint8(),
		Rejoin:         dataOut.ReadUint8(),
	}, nil
}

/*
	ZDO_MGMT_LQI_RSP

	This callback message is in response to the ZDO Management LQI Request.

	Usage:
		AREQ:
			         1         |      1      |       1     |    2    |    1   |            1         |      1     |             1          |            0-66
			Length = 0x06-0x48 | Cmd0 = 0x45 | Cmd1 = 0xB1 | SrcAddr | Status | NeighborTableEntries | StartIndex | NeighborTableListCount | NeighborTableListRecords
		Attributes:
			SrcAddr              2 bytes    Source address of the message
			Status               1 byte     This field indicates either SUCCESS or FAILURE.
			NeighborTableEntries 1 byte     Total number of entries available in the device.
			StartIndex           1 byte     Where in the total number of entries this response starts.
			NeighborLqiListCount 1 byte     Number of entries in this response.
			NeighborLqiList      0-66 bytes An array of NeighborLqiList items. NeighborLQICount contains the number of items in this table.
			                                ExtendedPanID                          8 bytes
			                                ExtendedAddress                        8 bytes
			                                NetworkAddress                         2 bytes
			                                DeviceType/ RxOnWhenIdle/ Relationship 1 byte
			                                PermitJoining                          1 byte
			                                Depth                                  1 byte
			                                LQI                                    1 byte
*/
func (c *Client) ZDOLQIMessage(frame *Frame) (*ZDOLQIMessage, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != 0xB1 {
		return nil, errors.New("frame isn't a LQI message")
	}

	dataOut := frame.DataAsBuffer()

	msg := &ZDOLQIMessage{
		SourceAddress:        dataOut.ReadUint16(),
		Status:               dataOut.ReadUint8(),
		NeighborTableEntries: dataOut.ReadUint8(),
		StartIndex:           dataOut.ReadUint8(),
		NeighborLQIListCount: dataOut.ReadUint8(),
		NeighborLqiList:      make([]NeighborLqiListItem, 0, 3),
	}

	if msg.Status != 0 {
		return nil, errors.New("failure")
	}

	for i := uint8(0); i < msg.NeighborLQIListCount; i++ {
		item := NeighborLqiListItem{
			ExtendedPanID:   dataOut.ReadIEEEAddr(),
			ExtendedAddress: dataOut.ReadIEEEAddr(),
			NetworkAddress:  dataOut.ReadUint16(),
		}

		v := dataOut.ReadUint8()
		item.DeviceType = v & 0x03
		item.RxOnWhenIdle = (v & 0x0C) >> 2
		item.Relationship = (v & 0x70) >> 4
		item.PermitJoining = dataOut.ReadUint8() & 0x03
		item.Depth = dataOut.ReadUint8()
		item.LQI = dataOut.ReadUint8()

		msg.NeighborLqiList = append(msg.NeighborLqiList, item)
	}

	return msg, nil
}

func (c *Client) ZDONetworkDiscoveryMessage(frame *Frame) (*ZDOManagementNetworkDiscoveryMessage, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != CommandManagementNetworkDiscoveryResponse {
		return nil, errors.New("frame isn't a network discovery message")
	}

	dataOut := frame.DataAsBuffer()

	msg := &ZDOManagementNetworkDiscoveryMessage{
		SourceAddress:    dataOut.ReadUint16(),
		Status:           dataOut.ReadUint8(),
		NetworkCount:     dataOut.ReadUint8(),
		StartIndex:       dataOut.ReadUint8(),
		NetworkListCount: dataOut.ReadUint8(),
		NetworkList:      make([]NetworkListItem, 0, 12),
	}

	if msg.Status != 0 {
		return nil, errors.New("failure")
	}

	for i := uint8(0); i < msg.NetworkListCount; i++ {
		item := NetworkListItem{
			PAN:            dataOut.ReadUint16(),
			LogicalChannel: dataOut.ReadUint8(),
		}

		v := dataOut.ReadUint8()
		item.StackProfile = v & 0x0F
		item.ZigBeeVersion = (v & 0xF0) >> 4

		v = dataOut.ReadUint8()
		item.BeaconOrder = v & 0x0F
		item.SuperFrameOrder = (v & 0xF0) >> 4

		msg.NetworkList = append(msg.NetworkList, item)
	}

	return msg, nil
}

func (c *Client) ZDOManagementRoutingTableMessage(frame *Frame) (*ZDOManagementNetworkDiscoveryMessage, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != CommandManagementRoutingTableResponse {
		return nil, errors.New("frame isn't a routing table message")
	}

	dataOut := frame.DataAsBuffer()

	fmt.Println("TEST", dataOut.Bytes())

	msg := &ZDOManagementNetworkDiscoveryMessage{
		SourceAddress:    dataOut.ReadUint16(),
		Status:           dataOut.ReadUint8(),
		NetworkCount:     dataOut.ReadUint8(),
		StartIndex:       dataOut.ReadUint8(),
		NetworkListCount: dataOut.ReadUint8(),
		NetworkList:      make([]NetworkListItem, 0, 12),
	}

	if msg.Status != 0 {
		return nil, errors.New("failure")
	}

	for i := uint8(0); i < msg.NetworkListCount; i++ {
		item := NetworkListItem{
			PAN:            dataOut.ReadUint16(),
			LogicalChannel: dataOut.ReadUint8(),
		}

		v := dataOut.ReadUint8()
		item.StackProfile = v & 0x0F
		item.ZigBeeVersion = (v & 0xF0) >> 4

		v = dataOut.ReadUint8()
		item.BeaconOrder = v & 0x0F
		item.SuperFrameOrder = (v & 0xF0) >> 4

		msg.NetworkList = append(msg.NetworkList, item)
	}

	return msg, nil
}
