package z_stack

import (
	"errors"
)

type ZDOSimpleDescriptorResponse struct {
	SourceAddress  uint16
	Status         CommandStatus
	NetworkAddress uint16
	Length         uint8
	Endpoint       uint8
	ProfileID      uint16
	DeviceID       uint16
	DeviceVersion  uint8
	NumInClusters  uint8
	InClusterList  []uint16
	NumOutClusters uint8
	OutClusterList []uint16
}

type ZDOActiveEndpointsResponse struct {
	SourceAddress  uint16
	Status         CommandStatus
	NetworkAddress uint16
	Count          uint8
	Endpoints      []uint8
}

type ZDODescriptionResponse struct {
	SourceAddress              uint16
	Status                     CommandStatus
	NetworkAddress             uint16
	LogicalType                uint8
	ComplexDescriptorAvailable uint8
	UserDescriptorAvailable    uint8
	APSFlags                   uint8
	FrequencyBand              uint8
	MacCapabilitiesFlags       uint8
	ManufacturerCode           uint16
	MaxBufferSize              uint8
	MaxInTransferSize          uint16
	ServerMask                 uint16
	MaxOutTransferSize         uint16
	DescriptorCapabilities     uint8
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

type ZDOLQIMessage struct {
	SourceAddress        uint16
	Status               CommandStatus
	NeighborTableEntries uint8
	StartIndex           uint8
	NeighborLQIListCount uint8
	NeighborLqiList      []NeighborLqiListItem
}

type ZDOManagementNetworkDiscoveryMessage struct {
	SourceAddress    uint16
	Status           CommandStatus
	NetworkCount     uint8
	StartIndex       uint8
	NetworkListCount uint8
	NetworkList      []NetworkListItem
}

type ZDOManagementRoutingTableMessage struct {
	SourceAddress         uint16
	Status                CommandStatus
	RoutingTableEntries   uint8
	StartIndex            uint8
	RoutingTableListCount uint8
	RoutingTableList      []RoutingTableListItem
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

type NetworkListItem struct {
	PAN             uint16
	LogicalChannel  uint8
	StackProfile    uint8
	ZigBeeVersion   uint8
	BeaconOrder     uint8
	SuperFrameOrder uint8
	PermitJoining   bool
}

type RoutingTableListItem struct {
	DestinationAddress uint16
	Status             uint8
	NextHop            uint16
}

/*
	ZDO_SIMPLE_DESC_RSP

	This callback message is in response to the ZDO Simple Descriptor Request

	Usage:
		AREQ:
			Length = 0x06-4E | Cmd0 = 0x45 | Cmd1 = 0x84 | SrcAddr | Status | NwkAddr | Len | Endpoint | ProfileId | DeviceId | DeviceVersion | NumInClusters | InClusterList | NumOutClusters | OutClusterList
		Attributes:
			SrcAddr        2 bytes    Specifies the message’s source network address.
			Status         1 byte     This field indicates either SUCCESS or FAILURE.
			NWKAddr        2 bytes    Specifies Device’s short address that this response describes.
			Len            1 byte     Specifies the length of the simple descriptor
			Endpoint       1 byte     Specifies Endpoint of the device
			ProfileId      2 bytes    The profile Id for this endpoint.
			DeviceId       2 bytes    The Device Description Id for this endpoint.
			DeviceVersion  1 byte     Defined as the following format
			                          0 – Version 1.00
			                          0x01-0x0F – Reserved.
			NumInClusters  1 byte     The number of input clusters in the InClusterList.
			InClusterList  0-32 bytes List of input cluster Id’s supported.
			NumOutClusters 1 byte     The number of output clusters in the OutClusterList.
			OutClusterList 0-32 bytes List of output cluster Id’s supported.
*/
func ZDOSimpleDescriptorResponseParse(frame *Frame) (*ZDOSimpleDescriptorResponse, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != CommandSimpleDescriptorResponse {
		return nil, errors.New("frame isn't a simple descriptor response")
	}

	dataOut := frame.DataAsBuffer()

	msg := &ZDOSimpleDescriptorResponse{
		SourceAddress: dataOut.ReadUint16(),
		Status:        dataOut.ReadCommandStatus(),
	}

	if msg.Status != CommandStatusSuccess {
		return nil, msg.Status
	}

	msg.NetworkAddress = dataOut.ReadUint16()
	msg.Length = dataOut.ReadUint8()
	msg.Endpoint = dataOut.ReadUint8()
	msg.ProfileID = dataOut.ReadUint16()
	msg.DeviceID = dataOut.ReadUint16()
	msg.DeviceVersion = dataOut.ReadUint8()

	msg.NumInClusters = dataOut.ReadUint8()
	msg.InClusterList = make([]uint16, 0, msg.NumInClusters)
	for i := uint8(1); i <= msg.NumInClusters; i++ {
		msg.InClusterList = append(msg.InClusterList, dataOut.ReadUint16())
	}

	msg.NumOutClusters = dataOut.ReadUint8()
	msg.OutClusterList = make([]uint16, 0, msg.NumOutClusters)
	for i := uint8(1); i <= msg.NumOutClusters; i++ {
		msg.OutClusterList = append(msg.OutClusterList, dataOut.ReadUint16())
	}

	return msg, nil
}

/*
	ZDO_ACTIVE_EP_RSP

	This callback message is in response to the ZDO Active Endpoint Request.

	Usage:
		AREQ:
			          1        |      1      |      1      |    2    |    1   |    2    |       1       |    0-77
			Length = 0x06-0x53 | Cmd0 = 0x45 | Cmd1 = 0x85 | SrcAddr | Status | NwkAddr | ActiveEPCount | ActiveEPList
		Attributes:
			SrcAddr       2 bytes    The message’s source network address.
			Status        1 byte     This field indicates either SUCCESS or FAILURE.
			NWKAddr       2 bytes    Device’s short address that this response describes.
			ActiveEPCount 1 byte     Number of active endpoint in the list
			ActiveEPList  0-77 bytes Array of active endpoints on this device.
*/
func ZDOActiveEndpointsResponseParse(frame *Frame) (*ZDOActiveEndpointsResponse, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != CommandActiveEndpointsResponse {
		return nil, errors.New("frame isn't a active endpoint response")
	}

	dataOut := frame.DataAsBuffer()

	msg := &ZDOActiveEndpointsResponse{
		SourceAddress: dataOut.ReadUint16(),
		Status:        dataOut.ReadCommandStatus(),
	}

	if msg.Status != CommandStatusSuccess {
		return nil, msg.Status
	}

	msg.NetworkAddress = dataOut.ReadUint16()
	msg.Count = dataOut.ReadUint8()
	msg.Endpoints = dataOut.Bytes()

	return msg, nil
}

/*
	ZDO_NODE_DESC_RSP

	This callback message is in response to the ZDO Node Descriptor Request.

	Usage:
		AREQ:
			       1      |      1      |      1      |    2    |    1   |    2    |           1           |        1      |          1         |         2        |       1       |        2        |      2     |          2         |           1
			Length = 0x12 | Cmd0 = 0x45 | Cmd1 = 0x82 | SrcAddr | Status | NwkAddr | LogicalType/          | APSFlags/     | MACCapabilityFlags | ManufacturerCode | MaxBufferSize | MaxTransferSize | ServerMask | MaxOutTransferSize | DescriptorCapabilities
			              |             |             |         |        |         | ComplexDescAvailable/ | FrequencyBand |                    |                  |               |                 |            |                    |
			              |             |             |         |        |         | UserDescAvailable/    |               |                    |                  |               |                 |            |                    |
		Attributes:
			SrcAddr                     2 bytes The message’s source network address.
			Status                      1 byte  This field indicates either SUCCESS or FAILURE.
			NWKAddrOfInterest           2 bytes Device’s short address of this Node descriptor
			LogicalType/                1 byte  Logical Type: Bit 0-2
			ComplexDescriptorAvailable/         0 - ZigBee Coordinator
			UserDescriptorAvailable             1 - ZigBee Router
			                                    2 - ZigBee End Device
			                                    ComplexDescriptorAvailable: Bit 4 – Indicates if complex descriptor is available for the node
			                                    NodeFrequencyBand – Bit 5-7 – Identifies node frequency band capabilities
			APSFlags/                   1 byte  APSFlags – Bit 0-4 – Node Flags assigned for APS. For V1.0 all bits are reserved.
			FrequencyBand                       NodeFrequencyBand – Bit 5-7 – Identifies node frequency band capabilities
			MacCapabilitiesFlags        1 byte  Capability flags stored for the MAC
			                                    0x00 - CAPINFO_DEVICETYPE_RFD
			                                    0x01 - CAPINFO_ALTPANCOORD
			                                    0x02 - CAPINFO_DEVICETYPE_FFD
			                                    0x04 - CAPINFO_POWER_AC
			                                    0x08 - CAPINFO_RCVR_ON_IDLE
			                                    0x40 - CAPINFO_SECURITY_CAPABLE
			                                    0x80 - CAPINFO_ALLOC_ADDR
			ManufacturerCode            2 bytes Specifies a manufacturer code that is allocated by the ZigBee Alliance, relating to the manufacturer to the device.
			MaxBufferSize               1 byte  Indicates size of maximum NPDU. This field is used as a high level indication for management.
			MaxInTransferSize           2 bytes Indicates maximum size of Transfer up to 0x7fff (This field is reserved in version 1.0 and shall be set to zero).
			ServerMask                  2 bytes Bit 0 - Primary Trust Center
			                                    1 - Backup Trust Center
			                                    2 - Primary Binding Table Cache
			                                    3 - Backup Binding Table Cache
			                                    4 - Primary Discovery Cache
			                                    5 - Backup Discovery Cache
			MaxOutTransferSize          2 bytes Indicates maximum size of Transfer up to 0x7fff (This field is reserved in version 1.0 and shall be set to zero).
			DescriptorCapabilities      1 byte  Specifies the Descriptor capabilities
*/
func ZDONodeDescriptionResponseParse(frame *Frame) (*ZDODescriptionResponse, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != CommandNodeDescriptionResponse {
		return nil, errors.New("frame isn't a node description response")
	}

	dataOut := frame.DataAsBuffer()

	msg := &ZDODescriptionResponse{
		SourceAddress: dataOut.ReadUint16(),
		Status:        dataOut.ReadCommandStatus(),
	}

	if msg.Status != CommandStatusSuccess {
		return nil, msg.Status
	}

	msg.NetworkAddress = dataOut.ReadUint16()

	b := dataOut.ReadUint8()
	msg.LogicalType = b & 0x07
	msg.ComplexDescriptorAvailable = b
	msg.UserDescriptorAvailable = b

	b = dataOut.ReadUint8()
	msg.APSFlags = b
	msg.FrequencyBand = msg.APSFlags

	msg.MacCapabilitiesFlags = dataOut.ReadUint8()
	msg.ManufacturerCode = dataOut.ReadUint16()
	msg.MaxBufferSize = dataOut.ReadUint8()
	msg.MaxInTransferSize = dataOut.ReadUint16()
	msg.ServerMask = dataOut.ReadUint16()
	msg.MaxOutTransferSize = dataOut.ReadUint16()
	msg.DescriptorCapabilities = dataOut.ReadUint8()

	return msg, nil
}

func ZDODeviceJoinedMessageParse(frame *Frame) (*ZDODeviceJoinedMessage, error) {
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
func ZDOEndDeviceAnnounceMessageParse(frame *Frame) (*ZDOEndDeviceAnnounceMessage, error) {
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

func ZDODeviceLeaveMessageParse(frame *Frame) (*ZDODeviceLeaveMessage, error) {
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
func ZDOLQIMessageParse(frame *Frame) (*ZDOLQIMessage, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != 0xB1 {
		return nil, errors.New("frame isn't a LQI message")
	}

	dataOut := frame.DataAsBuffer()
	if dataOut.Len() == 0 {
		return nil, errors.New("failure")
	}

	msg := &ZDOLQIMessage{
		SourceAddress:        dataOut.ReadUint16(),
		Status:               dataOut.ReadCommandStatus(),
		NeighborTableEntries: dataOut.ReadUint8(),
		StartIndex:           dataOut.ReadUint8(),
		NeighborLQIListCount: dataOut.ReadUint8(),
		NeighborLqiList:      make([]NeighborLqiListItem, 0, 3),
	}

	if msg.Status != CommandStatusSuccess {
		return nil, msg.Status
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

func ZDONetworkDiscoveryMessageParse(frame *Frame) (*ZDOManagementNetworkDiscoveryMessage, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != CommandManagementNetworkDiscoveryResponse {
		return nil, errors.New("frame isn't a network discovery message")
	}

	dataOut := frame.DataAsBuffer()

	msg := &ZDOManagementNetworkDiscoveryMessage{
		SourceAddress:    dataOut.ReadUint16(),
		Status:           dataOut.ReadCommandStatus(),
		NetworkCount:     dataOut.ReadUint8(),
		StartIndex:       dataOut.ReadUint8(),
		NetworkListCount: dataOut.ReadUint8(),
		NetworkList:      make([]NetworkListItem, 0, 12),
	}

	if msg.Status != CommandStatusSuccess {
		return nil, msg.Status
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

func ZDOManagementRoutingTableMessageParse(frame *Frame) (*ZDOManagementRoutingTableMessage, error) {
	if frame.SubSystem() != SubSystemZDOInterface {
		return nil, errors.New("frame isn't a ZDO interface")
	}

	if frame.CommandID() != CommandManagementRoutingTableResponse {
		return nil, errors.New("frame isn't a routing table message")
	}

	dataOut := frame.DataAsBuffer()

	msg := &ZDOManagementRoutingTableMessage{
		SourceAddress:         dataOut.ReadUint16(),
		Status:                dataOut.ReadCommandStatus(),
		RoutingTableEntries:   dataOut.ReadUint8(),
		StartIndex:            dataOut.ReadUint8(),
		RoutingTableListCount: dataOut.ReadUint8(),
		RoutingTableList:      make([]RoutingTableListItem, 0, 15),
	}

	if msg.Status != CommandStatusSuccess {
		return nil, msg.Status
	}

	for i := uint8(0); i < msg.RoutingTableListCount; i++ {
		msg.RoutingTableList = append(msg.RoutingTableList, RoutingTableListItem{
			DestinationAddress: dataOut.ReadUint16(),
			Status:             dataOut.ReadUint8(),
			NextHop:            dataOut.ReadUint16(),
		})
	}

	return msg, nil
}
