package devices

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ghthor/gowol"
	"github.com/gorilla/websocket"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/snabb/webostv"
)

const (
	TVLGWebOSMQTTTopicApplication        mqtt.Topic = boggart.ComponentName + "/tv/+/application"
	TVLGWebOSMQTTTopicMute               mqtt.Topic = boggart.ComponentName + "/tv/+/mute"
	TVLGWebOSMQTTTopicVolume             mqtt.Topic = boggart.ComponentName + "/tv/+/volume"
	TVLGWebOSMQTTTopicVolumeUp           mqtt.Topic = boggart.ComponentName + "/tv/+/volume/up"
	TVLGWebOSMQTTTopicVolumeDown         mqtt.Topic = boggart.ComponentName + "/tv/+/volume/down"
	TVLGWebOSMQTTTopicToast              mqtt.Topic = boggart.ComponentName + "/tv/+/toast"
	TVLGWebOSMQTTTopicPower              mqtt.Topic = boggart.ComponentName + "/tv/+/power"
	TVLGWebOSMQTTTopicStateMute          mqtt.Topic = boggart.ComponentName + "/tv/+/state/mute"
	TVLGWebOSMQTTTopicStateVolume        mqtt.Topic = boggart.ComponentName + "/tv/+/state/volume"
	TVLGWebOSMQTTTopicStateApplication   mqtt.Topic = boggart.ComponentName + "/tv/+/state/application"
	TVLGWebOSMQTTTopicStateChannelNumber mqtt.Topic = boggart.ComponentName + "/tv/+/state/channel-number"
	TVLGWebOSMQTTTopicStatePower         mqtt.Topic = boggart.ComponentName + "/tv/+/state/power"
)

var defaultDialerLGWebOSTV = webostv.Dialer{
	DisableTLS: true,
	WebsocketDialer: &websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		NetDial: (&net.Dialer{
			Timeout:   time.Second * 5,
			KeepAlive: time.Second * 30, // ensure we notice if the TV goes away
		}).Dial,
	},
}

type LGWebOSTV struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	mutex  sync.RWMutex
	client *webostv.Tv
	host   string
	key    string
}

func NewLGWebOSTV(host, key string) *LGWebOSTV {
	device := &LGWebOSTV{
		host: host,
		key:  key,
	}
	device.Init()
	device.SetDescription("LG TV WebOS")

	return device
}

func (d *LGWebOSTV) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeTV,
	}
}

func (d *LGWebOSTV) Ping(ctx context.Context) bool {
	client, err := d.Client()
	if err != nil {
		return false
	}

	_, err = client.GetCurrentSWInformation()
	return err == nil
}

func (d *LGWebOSTV) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillStopTask(d.taskSerialNumber)
	taskSerialNumber.SetTimeout(time.Second * 30)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Minute)
	taskSerialNumber.SetName("device-tv-lg-webos-serial-number")

	return []workers.Task{
		taskSerialNumber,
	}
}

func (d *LGWebOSTV) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		TVLGWebOSMQTTTopicMute,
		TVLGWebOSMQTTTopicVolume,
		TVLGWebOSMQTTTopicVolumeUp,
		TVLGWebOSMQTTTopicVolumeDown,
		TVLGWebOSMQTTTopicToast,
		TVLGWebOSMQTTTopicPower,
		TVLGWebOSMQTTTopicStateMute,
		TVLGWebOSMQTTTopicStateVolume,
		TVLGWebOSMQTTTopicStateApplication,
		TVLGWebOSMQTTTopicStateChannelNumber,
		TVLGWebOSMQTTTopicStatePower,
	}
}

func (d *LGWebOSTV) Client() (*webostv.Tv, error) {
	d.mutex.RLock()
	c := d.client
	d.mutex.RUnlock()

	if c != nil {
		return c, nil
	}

	client, err := defaultDialerLGWebOSTV.Dial(d.host)
	if err != nil {
		return nil, err
	}

	go client.MessageHandler()

	d.mutex.Lock()
	d.client = client
	d.mutex.Unlock()

	return client, nil
}

func (d *LGWebOSTV) taskSerialNumber(ctx context.Context) (interface{}, error, bool) {
	//if !d.IsEnabled() {
	//	return nil, nil, false
	//}

	client, err := d.Client()
	if err != nil {
		return nil, err, false
	}

	_, err = client.Register(d.key)
	if err != nil {
		return nil, err, false
	}

	if d.SerialNumber() == "" {
		deviceInfo, err := client.GetCurrentSWInformation()
		if err != nil {
			return nil, err, false
		}

		d.SetSerialNumber(deviceInfo.DeviceId)
	}

	sn := strings.Replace(d.SerialNumber(), ":", "-", -1)

	// set tv subscribers
	quit := make(chan struct{})

	go func() {
		state, err := client.ApplicationManagerGetForegroundAppInfo()
		if err == nil {
			d.monitorForegroundAppInfo(state)
		}

		client.ApplicationManagerMonitorForegroundAppInfo(d.monitorForegroundAppInfo, quit)
	}()

	go func() {
		state, err := client.AudioGetStatus()
		if err == nil {
			d.monitorAudio(state)
		}

		client.AudioMonitorStatus(d.monitorAudio, quit)
	}()

	go func() {
		state, err := client.TvGetCurrentChannel()
		if err == nil {
			d.monitorTvCurrentChannel(state)
		}

		client.TvMonitorCurrentChannel(d.monitorTvCurrentChannel, quit)
	}()

	// set MQTT subscribers
	err = d.MQTTSubscribe(TVLGWebOSMQTTTopicApplication.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		client.ApplicationManagerLaunch(string(message.Payload()), nil)
	})
	if err != nil {
		return nil, err, false
	}

	err = d.MQTTSubscribe(TVLGWebOSMQTTTopicMute.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		client.AudioSetMute(bytes.Equal(message.Payload(), []byte(`1`)))
	})
	if err != nil {
		return nil, err, false
	}

	err = d.MQTTSubscribe(TVLGWebOSMQTTTopicVolume.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		vol, err := strconv.ParseInt(string(message.Payload()), 10, 64)
		if err == nil {
			client.AudioSetVolume(int(vol))
		}
	})
	if err != nil {
		return nil, err, false
	}

	err = d.MQTTSubscribe(TVLGWebOSMQTTTopicVolumeUp.Format(sn), 0, func(_ context.Context, _ mqtt.Component, _ mqtt.Message) {
		client.AudioVolumeUp()
	})
	if err != nil {
		return nil, err, false
	}

	err = d.MQTTSubscribe(TVLGWebOSMQTTTopicVolumeDown.Format(sn), 0, func(_ context.Context, _ mqtt.Component, _ mqtt.Message) {
		client.AudioVolumeDown()
	})
	if err != nil {
		return nil, err, false
	}

	err = d.MQTTSubscribe(TVLGWebOSMQTTTopicToast.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		client.SystemNotificationsCreateToast(string(message.Payload()))
	})
	if err != nil {
		return nil, err, false
	}

	err = d.MQTTSubscribe(TVLGWebOSMQTTTopicPower.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		if bytes.Equal(message.Payload(), []byte(`1`)) {
			wol.MagicWake(d.SerialNumber(), "255.255.255.255")
		} else {
			client.SystemTurnOff()
		}
	})
	if err != nil {
		return nil, err, false
	}

	return nil, nil, true
}

func (d *LGWebOSTV) monitorForegroundAppInfo(s webostv.ForegroundAppInfo) error {
	ctx := context.Background()
	sn := strings.Replace(d.SerialNumber(), ":", "-", -1)

	if err := d.MQTTPublish(ctx, TVLGWebOSMQTTTopicStateApplication.Format(sn), 2, true, s.AppId); err != nil {
		return err
	}

	// TODO: cache
	if s.AppId == "" {
		if err := d.MQTTPublish(ctx, TVLGWebOSMQTTTopicStatePower.Format(sn), 2, true, false); err != nil {
			return err
		}
	} else {
		if err := d.MQTTPublish(ctx, TVLGWebOSMQTTTopicStatePower.Format(sn), 2, true, true); err != nil {
			return err
		}
	}

	return nil
}

func (d *LGWebOSTV) monitorAudio(s webostv.AudioStatus) error {
	ctx := context.Background()
	sn := strings.Replace(d.SerialNumber(), ":", "-", -1)

	if err := d.MQTTPublish(ctx, TVLGWebOSMQTTTopicStateMute.Format(sn), 2, true, s.Mute); err != nil {
		return err
	}

	if err := d.MQTTPublish(ctx, TVLGWebOSMQTTTopicStateVolume.Format(sn), 2, true, s.Volume); err != nil {
		return err
	}

	return nil
}

func (d *LGWebOSTV) monitorTvCurrentChannel(s webostv.TvCurrentChannel) error {
	ctx := context.Background()
	sn := strings.Replace(d.SerialNumber(), ":", "-", -1)

	if err := d.MQTTPublish(ctx, TVLGWebOSMQTTTopicStateChannelNumber.Format(sn), 2, true, s.ChannelNumber); err != nil {
		return err
	}

	return nil
}
