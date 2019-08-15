package xmeye

type CmdResponse struct {
	Name      string
	Ret       uint64
	SessionID string
}

type LoginResponse struct {
	AliveInterval uint64
	ChannelNum    uint64
	DataUseAES    bool
	DeviceType    string `json:"DeviceType "`
	ExtraChannel  uint64
	Ret           uint64
	SessionID     string
}

type OPTimeQuery struct {
	CmdResponse
	OPTimeQuery string
}

type SystemInfo struct {
	CmdResponse
	SystemInfo SystemInfoDetails
}

type SystemInfoDetails struct {
	AlarmInChannel  uint64
	AlarmOutChannel uint64
	AudioInChannel  uint64
	BuildTime       string
	CombineSwitch   uint64
	DeviceRunTime   string
	DigChannel      uint64
	EncryptVersion  string
	ExtraChannel    uint64
	HardWare        string
	HardWareVersion string
	SerialNo        string
	SoftWareVersion string
	TalkInChannel   uint64
	TalkOutChannel  uint64
	UpdataTime      string
	UpdataType      string
	VideoInChannel  uint64
	VideoOutChannel uint64
}
