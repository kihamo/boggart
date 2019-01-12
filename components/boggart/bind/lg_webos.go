package bind

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"strconv"
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
	LGWebOSMQTTTopicApplication        mqtt.Topic = boggart.ComponentName + "/tv/+/application"
	LGWebOSMQTTTopicMute               mqtt.Topic = boggart.ComponentName + "/tv/+/mute"
	LGWebOSMQTTTopicVolume             mqtt.Topic = boggart.ComponentName + "/tv/+/volume"
	LGWebOSMQTTTopicVolumeUp           mqtt.Topic = boggart.ComponentName + "/tv/+/volume/up"
	LGWebOSMQTTTopicVolumeDown         mqtt.Topic = boggart.ComponentName + "/tv/+/volume/down"
	LGWebOSMQTTTopicToast              mqtt.Topic = boggart.ComponentName + "/tv/+/toast"
	LGWebOSMQTTTopicPower              mqtt.Topic = boggart.ComponentName + "/tv/+/power"
	LGWebOSMQTTTopicStateMute          mqtt.Topic = boggart.ComponentName + "/tv/+/state/mute"
	LGWebOSMQTTTopicStateVolume        mqtt.Topic = boggart.ComponentName + "/tv/+/state/volume"
	LGWebOSMQTTTopicStateApplication   mqtt.Topic = boggart.ComponentName + "/tv/+/state/application"
	LGWebOSMQTTTopicStateChannelNumber mqtt.Topic = boggart.ComponentName + "/tv/+/state/channel-number"
	LGWebOSMQTTTopicStatePower         mqtt.Topic = boggart.ComponentName + "/tv/+/state/power"
)

var defaultDialerLGWebOS = webostv.Dialer{
	DisableTLS: true,
	WebsocketDialer: &websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		NetDial: (&net.Dialer{
			Timeout:   time.Second * 5,
			KeepAlive: time.Second * 30, // ensure we notice if the TV goes away
		}).Dial,
	},
}

type LGWebOS struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	mutex    sync.RWMutex
	initOnce sync.Once

	client *webostv.Tv
	config *LGWebOSConfig
}

type LGWebOSConfig struct {
	Host string `valid:"host,required"`
	Key  string `valid:"required"`
}

func (d LGWebOS) Config() interface{} {
	return &LGWebOSConfig{}
}

func (d LGWebOS) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	device := &LGWebOS{
		config: c.(*LGWebOSConfig),
	}
	device.Init()

	return device, nil
}

func (d *LGWebOS) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 10)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Second * 30)
	taskLiveness.SetName("bind-lg-webos-liveness")

	return []workers.Task{
		taskLiveness,
	}
}

func (d *LGWebOS) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		LGWebOSMQTTTopicMute,
		LGWebOSMQTTTopicVolume,
		LGWebOSMQTTTopicVolumeUp,
		LGWebOSMQTTTopicVolumeDown,
		LGWebOSMQTTTopicToast,
		LGWebOSMQTTTopicPower,
		LGWebOSMQTTTopicStateMute,
		LGWebOSMQTTTopicStateVolume,
		LGWebOSMQTTTopicStateApplication,
		LGWebOSMQTTTopicStateChannelNumber,
		LGWebOSMQTTTopicStatePower,
	}
}

func (d *LGWebOS) Client() (*webostv.Tv, error) {
	d.mutex.RLock()
	c := d.client
	d.mutex.RUnlock()

	if c != nil {
		return c, nil
	}

	client, err := defaultDialerLGWebOS.Dial(d.config.Host)
	if err != nil {
		return nil, err
	}

	go client.MessageHandler()

	d.mutex.Lock()
	d.client = client
	d.mutex.Unlock()

	return client, nil
}

func (d *LGWebOS) UpdateStatus(status boggart.DeviceStatus) {
	if status == boggart.DeviceStatusOffline && status != d.Status() {
		d.mutex.Lock()
		d.client = nil
		d.mutex.Unlock()
	}

	d.DeviceBindBase.UpdateStatus(status)
}

func (d *LGWebOS) taskLiveness(ctx context.Context) (interface{}, error) {
	client, err := d.Client()
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	_, err = client.Register(d.config.Key)
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	if d.Status() == boggart.DeviceStatusOnline {
		return nil, nil
	}

	if d.SerialNumber() == "" {
		deviceInfo, err := client.GetCurrentSWInformation()
		if err != nil {
			d.UpdateStatus(boggart.DeviceStatusOffline)
			return nil, err
		}

		d.SetSerialNumber(deviceInfo.DeviceId)
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)

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

func (d *LGWebOS) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(LGWebOSMQTTTopicApplication.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 2) {
				return
			}

			if client, err := d.Client(); err == nil {
				client.ApplicationManagerLaunch(string(message.Payload()), nil)
			}
		})),
		mqtt.NewSubscriber(LGWebOSMQTTTopicMute.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 2) {
				return
			}

			if client, err := d.Client(); err == nil {
				client.AudioSetMute(bytes.Equal(message.Payload(), []byte(`1`)))
			}
		})),
		mqtt.NewSubscriber(LGWebOSMQTTTopicVolume.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 2) {
				return
			}

			vol, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err == nil {
				if client, err := d.Client(); err == nil {
					client.AudioSetVolume(int(vol))
				}
			}
		})),
		mqtt.NewSubscriber(LGWebOSMQTTTopicVolumeUp.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 3) {
				return
			}

			if client, err := d.Client(); err == nil {
				client.AudioVolumeUp()
			}
		})),
		mqtt.NewSubscriber(LGWebOSMQTTTopicVolumeDown.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 3) {
				return
			}

			if client, err := d.Client(); err == nil {
				client.AudioVolumeDown()
			}
		})),
		mqtt.NewSubscriber(LGWebOSMQTTTopicToast.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 2) {
				return
			}

			if client, err := d.Client(); err == nil {
				client.SystemNotificationsCreateToast(string(message.Payload()))
			}
		})),
		mqtt.NewSubscriber(LGWebOSMQTTTopicPower.String(), 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) {
			if !boggart.CheckSerialNumberInMQTTTopic(d, message.Topic(), 2) {
				return
			}

			if bytes.Equal(message.Payload(), []byte(`1`)) {
				wol.MagicWake(d.SerialNumber(), "255.255.255.255")
			} else if d.Status() == boggart.DeviceStatusOnline {
				if client, err := d.Client(); err == nil {
					client.SystemTurnOff()
				}
			}
		}),
	}
}

func (d *LGWebOS) monitorForegroundAppInfo(s webostv.ForegroundAppInfo) error {
	ctx := context.Background()
	sn := mqtt.NameReplace(d.SerialNumber())

	d.MQTTPublishAsync(ctx, LGWebOSMQTTTopicStateApplication.Format(sn), 2, true, s.AppId)

	// TODO: cache
	if s.AppId == "" {
		d.UpdateStatus(boggart.DeviceStatusOffline)

		d.MQTTPublishAsync(ctx, LGWebOSMQTTTopicStatePower.Format(sn), 2, true, false)
	} else {
		d.MQTTPublishAsync(ctx, LGWebOSMQTTTopicStatePower.Format(sn), 2, true, true)
	}

	return nil
}

func (d *LGWebOS) monitorAudio(s webostv.AudioStatus) error {
	ctx := context.Background()
	sn := mqtt.NameReplace(d.SerialNumber())

	d.MQTTPublishAsync(ctx, LGWebOSMQTTTopicStateMute.Format(sn), 2, true, s.Mute)
	d.MQTTPublishAsync(ctx, LGWebOSMQTTTopicStateVolume.Format(sn), 2, true, s.Volume)

	return nil
}

func (d *LGWebOS) monitorTvCurrentChannel(s webostv.TvCurrentChannel) error {
	ctx := context.Background()
	sn := mqtt.NameReplace(d.SerialNumber())

	d.MQTTPublishAsync(ctx, LGWebOSMQTTTopicStateChannelNumber.Format(sn), 2, true, s.ChannelNumber)

	return nil
}
