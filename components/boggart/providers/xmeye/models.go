package xmeye

type Response struct {
	Name      string
	Ret       uint64
	SessionID string
}

type LoginResponse struct {
	AliveInterval uint64
	ChannelNum    uint64
	DataUseAES    bool
	DeviceType    string `json:DeviceType `
	ExtraChannel  uint64
	Ret           uint64
	SessionID     string
}

type OPTimeQuery struct {
	Response
	OPTimeQuery string
}

type SystemFunctions struct {
	Response
	SystemFunction struct {
		AlarmFunction struct {
			AlarmConfig           bool
			BlindDetect           bool
			HumanDVRSupportLevel  bool
			HumanDection          bool
			HumanDectionDVR       bool
			HumanDectionNVR       bool
			HumanDectionNVRNew    bool
			IPCAlarm              bool
			LossDetect            bool
			MotionDetect          bool
			NetAbort              bool
			NetAbortExtend        bool
			NetAlarm              bool
			NetIpConflict         bool
			NewVideoAnalyze       bool
			NewVideoAnalyze_digit bool
			PIRAlarm              bool
			SerialAlarm           bool
			StorageFailure        bool
			StorageLowSpace       bool
			StorageNotExist       bool
			VideoAnalyze          bool
		}
		CommFunction struct {
			CommRS232 bool
			CommRS485 bool
		}
		EncodeFunction struct {
			CombineStream      bool
			CustomChnDAMode    bool
			DoubleStream       bool
			IFrameRange        bool
			LowBitRate         bool
			MultiChannel       bool
			SmartEncodeDigital bool
			SmartH264          bool
			SmartH264V2        bool
			SnapStream         bool
			WaterMark          bool
		}
		InputMethod struct {
			NoSupportChinese bool
		}
		MobileDVR struct {
			CarPlateSet    bool
			DVRBootType    bool
			DelaySet       bool
			GpsTiming      bool
			StatusExchange bool
		}
		NetServerFunction struct {
			DualEthernet          bool
			IPAdaptive            bool
			MACProtocol           bool
			MonitorPlatform       bool
			NATProtocol           bool
			Net3G                 bool
			Net4G                 bool
			NetARSP               bool
			NetAlarmCenter        bool
			NetAnJuP2P            bool
			NetBaiduCloud         bool
			NetBjlThy             bool
			NetDAS                bool
			NetDDNS               bool
			NetDHCP               bool
			NetDNS                bool
			NetDataLink           bool
			NetEmail              bool
			NetEmailTLS           bool
			NetFTP                bool
			NetGodEyeAlarm        bool
			NetHMS                bool
			NetIPFilter           bool
			NetIPv6               bool
			NetKaiCong            bool
			NetKeyboard           bool
			NetLocalSdkPlatform   bool
			NetMidDAS             bool
			NetMobile             bool
			NetMobileWatch        bool
			NetMutliCast          bool
			NetNTP                bool
			NetNat                bool
			NetOpenVPN            bool
			NetPMS                bool
			NetPMSV2              bool
			NetPPPoE              bool
			NetPhoneMultimediaMsg bool
			NetPhoneShortMsg      bool
			NetPlatMega           bool
			NetPlatShiSou         bool
			NetPlatVVEye          bool
			NetPlatXingWang       bool
			NetRTSP               bool
			NetSPVMN              bool
			NetSPVMNSIP           bool
			NetTUTKIOTC           bool
			NetUPNP               bool
			NetVPN                bool
			NetWifi               bool
			NetWifiMode           bool
			OnvifPwdCheckout      bool
			PlatFormGBeyes        bool
			XMHeartBeat           bool
		}
		OtherFunction struct {
			AlterDigitalName           bool
			DownLoadPause              bool
			HddLowSpaceUseMB           bool
			HideDigital                bool
			MusicFilePlay              bool
			NOHDDRECORD                bool
			NotSupportAH               bool
			NotSupportAV               bool
			NotSupportTalk             bool
			SDsupportRecord            bool
			ShowAlarmLevelRegion       bool
			ShowFalseCheckTime         bool
			SupportAbnormitySendMail   bool
			SupportAdminContactInfo    bool
			SupportAlarmLinkLight      bool
			SupportAlarmVoiceTips      bool
			SupportAppBindFlag         bool
			SupportBT                  bool
			SupportC7Platform          bool
			SupportCamareStyle         bool
			SupportCameraMotorCtrl     bool
			SupportCameraWhiteLight    bool
			SupportCfgCloudupgrade     bool
			SupportCloudUpgrade        bool
			SupportCloudUpgradeIPC     bool
			SupportCoaxialParamCtrl    bool
			SupportCommDataUpload      bool
			SupportCorridorMode        bool
			SupportCustomOemInfo       bool
			SupportDeviceInfoNew       bool
			SupportDigitalEncode       bool
			SupportDigitalPre          bool
			SupportDimenCode           bool
			SupportDoubleLightBulb     bool
			SupportEncodeAddBeep       bool
			SupportExtreCode           bool
			SupportFTPTest             bool
			SupportFaceDetect          bool
			SupportFileUpgradeIPC      bool
			SupportFileUpgradeRouter   bool
			SupportFishEye             bool
			SupportImpRecord           bool
			SupportIntelligentPlayBack bool
			SupportInterruptSnap       bool
			SupportLimitNetLoginUsers  bool
			SupportLogStorageCtrl      bool
			SupportMailTest            bool
			SupportMaxPlayback         bool
			SupportModifyFrontcfg      bool
			SupportMusicLightBulb      bool
			SupportNVR                 bool
			SupportNetLocalSearch      bool
			SupportOSDInfo             bool
			SupportOnvifClient         bool
			SupportPOS                 bool
			SupportPWDSafety           bool
			SupportParkingGuide        bool
			SupportPlateDetect         bool
			SupportPlayBackExactSeek   bool
			SupportPlaybackLocate      bool
			SupportPtzIdleState        bool
			SupportRPSVideo            bool
			SupportRTSPClient          bool
			SupportResumePtzState      bool
			SupportSPVMNNasServer      bool
			SupportSafetyEmail         bool
			SupportSetDigIP            bool
			SupportSetRtcTime          bool
			SupportShowConnectStatus   bool
			SupportShowH265X           bool
			SupportShowProductType     bool
			SupportSmallChnTitleFont   bool
			SupportSnapCfg             bool
			SupportSnapSchedule        bool
			SupportSnapV2Stream        bool
			SupportSnapshotConfigV2    bool
			SupportSoftPhotosensitive  bool
			SupportSplitControl        bool
			SupportStatusLed           bool
			SupportStorageFailReboot   bool
			SupportStringChangedXPOE   bool
			SupportSwitchResolution    bool
			SupportTextPassword        bool
			SupportTimeSetNewWay       bool
			SupportTimeZone            bool
			SupportUserProgram         bool
			SupportWIFINVR             bool
			SupportWarnWeakPWD         bool
			SupportWifiHotSpot         bool
			SupportWriteLog            bool
			Supportonviftitle          bool
			SuppportChangeOnvifPort    bool
			TitleAndStateUpload        bool
			USBsupportRecord           bool
		}
		PreviewFunction struct {
			GUISet bool
			Tour   bool
		}
		TipShow struct {
			NoBeepTipShow  bool
			NoEmailTipShow bool
			NoFTPTipShow   bool
		}
	}
}

type SystemInfo struct {
	Response
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
	Response
	OEMInfo struct {
		Address   string
		Name      string
		OEMID     uint64
		Telephone string
	}
}

type StorageInfo struct {
	Response
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
	Response
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

type AlarmInfo struct {
	Response
	AlarmInfo struct {
		Channel   uint8
		Event     string
		StartTime string
		Status    string
	}
}
