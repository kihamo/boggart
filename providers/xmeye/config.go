package xmeye

import (
	"context"
	"io"
)

type configName string

const (
	ConfigNameAVEncEncode            configName = "AVEnc.Encode"
	ConfigNameAVEncVideoWidget       configName = "AVEnc.VideoWidget"
	ConfigNameAVEncVideoColor        configName = "AVEnc.VideoColor"
	ConfigNameAVEncCombineEncode     configName = "AVEnc.CombineEncode"
	ConfigNameDetectMotionDetect     configName = "Detect.MotionDetect"
	ConfigNameDetectBlindDetect      configName = "Detect.BlindDetect"
	ConfigNameDetectLossDetect       configName = "Detect.LossDetect"
	ConfigNameAlarmLocalAlarm        configName = "Alarm.LocalAlarm"
	ConfigNameAlarmNetAlarm          configName = "Alarm.NetAlarm"
	ConfigNameAlarmNetIPConflict     configName = "Alarm.NetIPConflict"
	ConfigNameAlarmNetAbort          configName = "Alarm.NetAbort"
	ConfigNameAlarmPTZAlarmProtocol  configName = "Alarm.PTZAlarmProtocol"
	ConfigNameStorageStorageNotExist configName = "Storage.StorageNotExist"
	ConfigNameStorageStorageLowSpace configName = "Storage.StorageLowSpace"
	ConfigNameStorageStorageFailure  configName = "Storage.StorageFailure"
	ConfigNameStorageSnapshot        configName = "Storage.Snapshot"
	ConfigNameNetworkNetCommon       configName = "NetWork.NetCommon"
	ConfigNameNetworkNetIPFilter     configName = "NetWork.NetIPFilter"
	ConfigNameNetworkNetDHCP         configName = "NetWork.NetDHCP"
	ConfigNameNetworkNetDDNS         configName = "NetWork.NetDDNS"
	ConfigNameNetworkNetEmail        configName = "NetWork.NetEmail"
	ConfigNameNetworkNetNTP          configName = "NetWork.NetNTP"
	ConfigNameNetworkNetPPPoE        configName = "NetWork.NetPPPoE"
	ConfigNameNetworkNetDNS          configName = "NetWork.NetDNS"
	ConfigNameNetworkNetARSP         configName = "NetWork.NetARSP"
	ConfigNameNetworkNetMobile       configName = "NetWork.NetMobile"
	ConfigNameNetworkUpnp            configName = "NetWork.Upnp"
	ConfigNameNetworkNetFTP          configName = "NetWork.NetFTP"
	ConfigNameNetworkAlarmServer     configName = "NetWork.AlarmServer"
	ConfigNameUartComm               configName = "Uart.Comm"
	ConfigNameUartPTZ                configName = "Uart.PTZ"
	ConfigNameUartPTZPreset          configName = "Uart.PTZPreset"
	ConfigNameUartPTZTour            configName = "Uart.PTZTour"
	ConfigNameFVideoTour             configName = "fVideo.Tour"
	ConfigNameFVideoGUISet           configName = "fVideo.GUISet"
	ConfigNameFVideoTVAdjust         configName = "fVideo.TVAdjust"
	ConfigNameFVideoAudioInFormat    configName = "fVideo.AudioInFormat"
	ConfigNameFVideoPlay             configName = "fVideo.Play"
	ConfigNameGeneralGeneral         configName = "General.General"
	ConfigNameGeneralLocation        configName = "General.Location"
	ConfigNameGeneralAutoMaintain    configName = "General.AutoMaintain"
	ConfigNameChannelTitle           configName = "ChannelTitle"
	ConfigNameRecord                 configName = "Record"
)

func (c *Client) ConfigGet(ctx context.Context, name configName, def bool) (interface{}, error) {
	var result map[string]interface{}

	code := CmdConfigGetRequest
	if def {
		code = CmdDefaultConfigGetRequest
	}

	err := c.CmdWithResult(ctx, code, string(name), &result)
	if err != nil {
		return nil, err
	}

	if values, ok := result[string(name)]; ok {
		return values, nil
	}

	return nil, err
}

func (c *Client) ConfigChannelTitleGet(ctx context.Context) ([]string, error) {
	var result struct {
		Response
		ChannelTitleGet []string
	}

	err := c.CmdWithResult(ctx, CmdConfigChannelTitleGetRequest, "ChannelTitleGet", &result)
	if err != nil {
		return nil, err
	}

	return result.ChannelTitleGet, err
}

func (c *Client) ConfigChannelTitleSet(ctx context.Context, names ...string) error {
	_, err := c.Call(ctx, CmdConfigChannelTitleSetRequest, map[string]interface{}{
		"Name":         "ChannelTitle",
		"ChannelTitle": names,
	})
	return err
}

func (c *Client) ConfigExport(ctx context.Context) (io.Reader, error) {
	packet, err := c.Call(ctx, CmdConfigExportRequest, nil)

	return packet.Payload, err
}
