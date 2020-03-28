package xmeye

import (
	"time"
)

type Response struct {
	Ret       uint64
	Name      string
	SessionID string
}

type LoginResponse struct {
	DataUseAES    bool
	SessionID     Uint32
	AliveInterval uint64
	ChannelNum    uint64
	ExtraChannel  uint64
	Ret           uint64
	DeviceType    string `json:"DeviceType "`
}

type LogSearch struct {
	Data     string
	Position uint32
	Time     Time
	Type     string
	User     string
}

type FileSearch struct {
	DiskNo     uint32
	SerialNo   uint32
	BeginTime  Time
	EndTime    Time
	FileLength Uint32
	FileName   string
}

func (f FileSearch) Duration() time.Duration {
	return f.EndTime.Sub(f.BeginTime.Time)
}

type SystemInfo struct {
	AlarmInChannel  uint64
	AlarmOutChannel uint64
	AudioInChannel  uint64
	CombineSwitch   uint64
	DigChannel      uint64
	ExtraChannel    uint64
	TalkInChannel   uint64
	TalkOutChannel  uint64
	VideoInChannel  uint64
	VideoOutChannel uint64
	BuildTime       Time
	UpdataTime      Time
	UpdataType      Uint32
	DeviceRunTime   Uint32
	EncryptVersion  string
	HardWare        string
	HardWareVersion string
	SerialNo        string
	SoftWareVersion string
}

type OEMInfo struct {
	Address   string
	Name      string
	OEMID     uint64
	Telephone string
}

type StorageInfo struct {
	PartNumber uint64
	PlysicalNo uint64
	Partition  []struct {
		DirverType    uint64
		IsCurrent     bool
		LogicSerialNo uint64
		NewEndTime    Time
		NewStartTime  Time
		OldEndTime    Time
		OldStartTime  Time
		Status        uint64
		RemainSpace   Uint32
		TotalSpace    Uint32
	}
}

type WorkState struct {
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

type AlarmInfo struct {
	Channel   uint8
	Event     string
	StartTime Time
	Status    string
}

type OPSystemUpgrade struct {
	Hardware string
	LogoArea struct {
		Begin Uint32
		End   Uint32
	}
	LogoPartType string
	Serial       string
	Vendor       string
}

type User struct {
	AuthorityList []string
	Group         string
	Memo          string
	Name          string
	NoMD5         string
	Password      string
	Reserved      bool
	Sharable      bool
}

type Group struct {
	AuthorityList []string
	Memo          string
	Name          string
}
