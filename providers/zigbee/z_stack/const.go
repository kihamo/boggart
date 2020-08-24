package zstack

//go:generate /bin/bash -c "enumer -type=DeviceState -trimprefix=DeviceState -output=device_state_enumer.go"
//go:generate /bin/bash -c "enumer -type=CommandStatus -trimprefix=CommandStatus -output=command_status_enumer.go"

import (
	"time"
)

type DeviceState uint8
type CommandStatus uint8

const (
	VersionZStack12  = 0
	VersionZStack3x0 = 1
	VersionZStack30x = 2

	// A POLL command is used to retrieve queued data. This command is only applicable to SPI transport.
	// For a POLL command the subsystem and Id are set to zero and data length is zero.
	TypePoll = 0x0

	// A synchronous request that requires an immediate response.
	// For example, a function call with a return value would use an SREQ command.
	TypeSREQ = 0x1

	// An asynchronous request.
	// For example, a callback event or a function call with no return value would use an AREQ command.
	TypeAREQ = 0x2

	// A synchronous response. This type of command is only sent in response to a SREQ command.
	// For an SRSP command the subsystem and Id are set to the same values as the corresponding SREQ.
	// The length of an SRSP is generally nonzero, so an SRSP with length=0 can be used to indicate an error.
	TypeSRSP = 0x3

	DeviceTypeNone        = 0
	DeviceTypeCoordinator = 1
	DeviceTypeRouter      = 2
	DeviceTypeEndDevice   = 4

	DeviceLogicalTypeCoordinator = 0
	DeviceLogicalTypeRouter      = 1
	DeviceLogicalTypeEndDevice   = 2

	SubSystemReserved       = 0x00
	SubSystemSysInterface   = 0x01
	SubSystemMACInterface   = 0x02
	SubSystemNWKInterface   = 0x03
	SubSystemAFInterface    = 0x04
	SubSystemZDOInterface   = 0x05
	SubSystemSAPIInterface  = 0x06
	SubSystemUtilInterface  = 0x07
	SubSystemDebugInterface = 0x08
	SubSystemAppInterface   = 0x09

	SOF            = byte(0xFE)
	FrameLengthMin = 5
	FrameLengthMax = 258

	PositionFrameLength = 1
	PositionCommand1    = 2
	PositionCommand2    = 3
	PositionData        = 4

	DeviceStateInitializedNotStarted   DeviceState = 0x00 // Initialized - not started automatically
	DeviceStateInitializedNotConnected DeviceState = 0x01 // Initialized - not connected to anything
	DeviceStateDiscoveringPAN          DeviceState = 0x02 // Discovering PAN's to join
	DeviceStateJoiningPAN              DeviceState = 0x03 // Joining a PAN
	DeviceStateRejoiningPAN            DeviceState = 0x04 // Rejoining a PAN, only for end devices
	DeviceStateUnauthentication        DeviceState = 0x05 // Joined but not yet authenticated by trust center
	DeviceStateStartedDevice           DeviceState = 0x06 // Started as device after authentication
	DeviceStateRouter                  DeviceState = 0x07 // Device joined, authenticated and is a router
	DeviceStateStartingCoordinator     DeviceState = 0x08 // Starting as ZigBee Coordinator
	DeviceStateStartedCoordinator      DeviceState = 0x09 // Started as ZigBee Coordinator
	DeviceStateOrphan                  DeviceState = 0x0A // Device has lost information about its parent

	ScanChannelsNone        = 0x00000000
	ScanChannelsAllChannels = 0x07FFF800
	ScanChannelsChannel11   = 0x00000800
	ScanChannelsChannel12   = 0x00001000
	ScanChannelsChannel13   = 0x00002000
	ScanChannelsChannel14   = 0x00004000
	ScanChannelsChannel15   = 0x00008000
	ScanChannelsChannel16   = 0x00010000
	ScanChannelsChannel17   = 0x00020000
	ScanChannelsChannel18   = 0x00040000
	ScanChannelsChannel19   = 0x00080000
	ScanChannelsChannel20   = 0x00100000
	ScanChannelsChannel21   = 0x00200000
	ScanChannelsChannel22   = 0x00400000
	ScanChannelsChannel23   = 0x00800000
	ScanChannelsChannel24   = 0x01000000
	ScanChannelsChannel25   = 0x02000000
	ScanChannelsChannel26   = 0x04000000

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

	CommandGetDeviceInfo                      = 0x00
	CommandLEDControl                         = 0x0A
	CommandAfIncomingMessage                  = 0x81
	CommandNodeDescriptionResponse            = 0x82
	CommandSimpleDescriptorResponse           = 0x84
	CommandActiveEndpointsResponse            = 0x85
	CommandManagementNetworkDiscoveryResponse = 0xB0
	CommandManagementRoutingTableResponse     = 0xB2
	CommandManagementPermitJoinResponse       = 0xB6
	CommandEndDeviceAnnounceInd               = 0xC1
	CommandLeaveInd                           = 0xC9
	CommandTcDeviceInd                        = 0xCA
	CommandPermitJoinInd                      = 0xCB

	ResetTypeHard = 0x0
	ResetTypeSoft = 0x1

	// https://github.com/Koenkk/zigbee-herdsman/blob/9299a5ffed13baa0007866ad10af8ef0d1bfb63d/src/adapter/z-stack/constants/common.ts#L159
	CommandStatusSuccess           CommandStatus = 0x00
	CommandStatusFailure           CommandStatus = 0x01
	CommandStatusInvalidParam      CommandStatus = 0x02
	CommandStatusNvItemInitialized CommandStatus = 0x09
	CommandStatusApsDuplicateEntry CommandStatus = 0xB8
	CommandStatusNwkInvalidRequest CommandStatus = 0xC2

	NvItemIDHasConfiguredZStack1 = 0x0F00
	NvItemIDHasConfiguredZStack3 = 0x0060

	ZdConfigurationNetworkKey = 0x62

	DefaultWaitTimeout = time.Millisecond * 6000

	InterviewStatusDefault = uint32(0) + 1
	InterviewStatusStarted
	InterviewStatusCompleted
)

func (i CommandStatus) Error() string {
	return i.String()
}
