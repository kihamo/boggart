package xmeye

import (
	"github.com/kihamo/boggart/components/boggart/providers/xmeye/internal"
)

const (
	ConfigNameAVEncEncode            = "AVEnc.Encode"
	ConfigNameAVEncVideoWidget       = "AVEnc.VideoWidget"
	ConfigNameAVEncVideoColor        = "AVEnc.VideoColor"
	ConfigNameAVEncCombineEncode     = "AVEnc.CombineEncode"
	ConfigNameDetectMotionDetect     = "Detect.MotionDetect"
	ConfigNameDetectBlindDetect      = "Detect.BlindDetect"
	ConfigNameDetectLossDetect       = "Detect.LossDetect"
	ConfigNameAlarmLocalAlarm        = "Alarm.LocalAlarm"
	ConfigNameAlarmNetAlarm          = "Alarm.NetAlarm"
	ConfigNameAlarmNetIPConflict     = "Alarm.NetIPConflict"
	ConfigNameAlarmNetAbort          = "Alarm.NetAbort"
	ConfigNameAlarmPTZAlarmProtocol  = "Alarm.PTZAlarmProtocol"
	ConfigNameStorageStorageNotExist = "Storage.StorageNotExist"
	ConfigNameStorageStorageLowSpace = "Storage.StorageLowSpace"
	ConfigNameStorageStorageFailure  = "Storage.StorageFailure"
	ConfigNameStorageSnapshot        = "Storage.Snapshot"
	ConfigNameNetWorkNetCommon       = "NetWork.NetCommon"
	ConfigNameNetWorkNetIPFilter     = "NetWork.NetIPFilter"
	ConfigNameNetWorkNetDHCP         = "NetWork.NetDHCP"
	ConfigNameNetWorkNetDDNS         = "NetWork.NetDDNS"
	ConfigNameNetWorkNetEmail        = "NetWork.NetEmail"
	ConfigNameNetWorkNetNTP          = "NetWork.NetNTP"
	ConfigNameNetWorkNetPPPoE        = "NetWork.NetPPPoE"
	ConfigNameNetWorkNetDNS          = "NetWork.NetDNS"
	ConfigNameNetWorkNetARSP         = "NetWork.NetARSP"
	ConfigNameNetWorkNetMobile       = "NetWork.NetMobile"
	ConfigNameNetWorkUpnp            = "NetWork.Upnp"
	ConfigNameNetWorkNetFTP          = "NetWork.NetFTP"
	ConfigNameNetWorkAlarmServer     = "NetWork.AlarmServer"
	ConfigNameUartComm               = "Uart.Comm"
	ConfigNameUartPTZ                = "Uart.PTZ"
	ConfigNameUartPTZPreset          = "Uart.PTZPreset"
	ConfigNameUartPTZTour            = "Uart.PTZTour"
	ConfigNameFVideoTour             = "fVideo.Tour"
	ConfigNameFVideoGUISet           = "fVideo.GUISet"
	ConfigNameFVideoTVAdjust         = "fVideo.TVAdjust"
	ConfigNameFVideoAudioInFormat    = "fVideo.AudioInFormat"
	ConfigNameFVideoPlay             = "fVideo.Play"
	ConfigNameGeneralGeneral         = "General.General"
	ConfigNameGeneralLocation        = "General.Location"
	ConfigNameGeneralAutoMaintain    = "General.AutoMaintain"
	ConfigNameChannelTitle           = "ChannelTitle"
	ConfigNameRecord                 = "Record"
)

func (c *Client) ConfigGet(name string, def bool) (map[string]interface{}, error) {
	var result map[string]interface{}

	code := CmdConfigGetRequest
	if def {
		code = CmdDefaultConfigGetRequest
	}

	err := c.CmdWithResult(code, name, &result)
	if err != nil {
		return nil, err
	}

	if values, ok := result[name]; ok {
		if config, ok := values.(map[string]interface{}); ok {
			return config, nil
		}
	}

	return nil, err
}

func (c *Client) ConfigChannelTitleGet() ([]string, error) {
	var result struct {
		internal.Response
		ChannelTitleGet []string
	}

	err := c.CmdWithResult(CmdConfigChannelTitleGetRequest, "ChannelTitleGet", &result)
	if err != nil {
		return nil, err
	}

	return result.ChannelTitleGet, err
}

func (c *Client) ConfigChannelTitleSet(names ...string) error {
	_, err := c.Call(CmdConfigChannelTitleSetRequest, map[string]interface{}{
		"Name":         "ChannelTitle",
		"ChannelTitle": names,
	})
	return err
}
