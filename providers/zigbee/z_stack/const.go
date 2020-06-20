package z_stack

const (
	// A POLL command is used to retrieve queued data. This command is only applicable to SPI transport.
	// For a POLL command the subsystem and Id are set to zero and data length is zero.
	TypePoll = 0

	// A synchronous request that requires an immediate response.
	// For example, a function call with a return value would use an SREQ command.
	TypeSREQ = 1

	// An asynchronous request.
	// For example, a callback event or a function call with no return value would use an AREQ command.
	TypeAREQ = 2

	// A synchronous response. This type of command is only sent in response to a SREQ command.
	// For an SRSP command the subsystem and Id are set to the same values as the corresponding SREQ.
	// The length of an SRSP is generally nonzero, so an SRSP with length=0 can be used to indicate an error.
	TypeSRSP = 3

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

	// Initialized - not started automatically 0x01: Initialized - not connected to anything 0x02: Discovering PAN's to join
	DeviceStateInitialized = 0x00

	// Joining a PAN
	DeviceStateJoining = 0x03

	// Rejoining a PAN, only for end devices
	DeviceStateRejoining = 0x04

	// Joined but not yet authenticated by trust center 0x06: Started as device after authentication
	DeviceStateUnauthentication = 0x05

	// Device joined, authenticated and is a router 0x08: Starting as ZigBee Coordinator
	DeviceStateRouter = 0x07

	// Started as ZigBee Coordinator
	DeviceStateCoordinator = 0x09

	// Device has lost information about its parent
	DeviceStateOrphan = 0x0A
)
