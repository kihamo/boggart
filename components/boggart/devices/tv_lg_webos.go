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

	mutex    sync.RWMutex
	initOnce sync.Once

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
	return true
}

func (d *LGWebOSTV) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 10)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Second * 30)
	taskLiveness.SetName("device-tv-lg-webos-liveness")

	return []workers.Task{
		taskLiveness,
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

func (d *LGWebOSTV) UpdateStatus(status boggart.DeviceStatus) {
	if status == boggart.DeviceStatusOffline && status != d.Status() {
		d.mutex.Lock()
		d.client = nil
		d.mutex.Unlock()
	}

	d.DeviceBase.UpdateStatus(status)
}

func (d *LGWebOSTV) taskLiveness(ctx context.Context) (interface{}, error) {
	client, err := d.Client()
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	_, err = client.Register(d.key)
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	if d.Status() == boggart.DeviceStatusOnline {
		return nil, nil
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)

	if d.SerialNumber() == "" {
		deviceInfo, err := client.GetCurrentSWInformation()
		if err != nil {
			d.UpdateStatus(boggart.DeviceStatusOffline)
			return nil, err
		}

		d.SetSerialNumber(deviceInfo.DeviceId)
		d.initOnce.Do(d.initMQTTSubscribers)
	}

	// set tv subscribers
	// TODO: close if OFFLINE
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

	return nil, nil
}

func (d *LGWebOSTV) initMQTTSubscribers() {
	sn := strings.Replace(d.SerialNumber(), ":", "-", -1)

	d.MQTTSubscribeDeviceIsOnline(TVLGWebOSMQTTTopicApplication.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		if client, err := d.Client(); err == nil {
			client.ApplicationManagerLaunch(string(message.Payload()), nil)
		}
	})

	d.MQTTSubscribeDeviceIsOnline(TVLGWebOSMQTTTopicMute.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		if client, err := d.Client(); err == nil {
			client.AudioSetMute(bytes.Equal(message.Payload(), []byte(`1`)))
		}
	})

	d.MQTTSubscribeDeviceIsOnline(TVLGWebOSMQTTTopicVolume.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		vol, err := strconv.ParseInt(string(message.Payload()), 10, 64)
		if err == nil {
			if client, err := d.Client(); err == nil {
				client.AudioSetVolume(int(vol))
			}
		}
	})

	d.MQTTSubscribeDeviceIsOnline(TVLGWebOSMQTTTopicVolumeUp.Format(sn), 0, func(_ context.Context, _ mqtt.Component, _ mqtt.Message) {
		if client, err := d.Client(); err == nil {
			client.AudioVolumeUp()
		}
	})

	d.MQTTSubscribeDeviceIsOnline(TVLGWebOSMQTTTopicVolumeDown.Format(sn), 0, func(_ context.Context, _ mqtt.Component, _ mqtt.Message) {
		if client, err := d.Client(); err == nil {
			client.AudioVolumeDown()
		}
	})

	d.MQTTSubscribeDeviceIsOnline(TVLGWebOSMQTTTopicToast.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		if client, err := d.Client(); err == nil {
			client.SystemNotificationsCreateToast(string(message.Payload()))
		}
	})

	d.MQTTSubscribe(TVLGWebOSMQTTTopicPower.Format(sn), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
		if bytes.Equal(message.Payload(), []byte(`1`)) {
			wol.MagicWake(d.SerialNumber(), "255.255.255.255")
		} else if d.Status() == boggart.DeviceStatusOnline {
			if client, err := d.Client(); err == nil {
				client.SystemTurnOff()
			}
		}
	})
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
