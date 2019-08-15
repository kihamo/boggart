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
	SystemInfo struct {
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
}

type OEMInfo struct {
	CmdResponse
	OEMInfo struct {
		Address   string
		Name      string
		OEMID     uint64
		Telephone string
	}
}

type StorageInfo struct {
	CmdResponse
	StorageInfo []struct {
		PartNumber uint64
		PlysicalNo uint64
		Partition  []struct {
			DirverType    uint64
			IsCurrent     bool
			LogicSerialNo uint64
			NewEndTime    string
			NewStartTime  string
			OldEndTime    string
			OldStartTime  string
			RemainSpace   string
			Status        uint64
			TotalSpace    string
		}
	}
}

type WorkState struct {
	CmdResponse
	WorkState struct {
		AlarmState struct {
			AlarmIn     uint64
			AlarmOut    uint64
			VideoBlind  uint64
			VideoLoss   uint64
			VideoMotion uint64
		}
		ChannelState []struct {
			Bitrate uint64
			Record  bool
		}
	}
}
