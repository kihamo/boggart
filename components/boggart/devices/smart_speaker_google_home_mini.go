package devices

import (
	"bytes"
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/google/home"
	"github.com/kihamo/boggart/components/boggart/providers/google/home/client"
	"github.com/kihamo/boggart/components/boggart/providers/google/home/client/info"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/voice/players/chromecast"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

const (
	GoogleHomeMiniSmartSpeakerMQTTTopicVolume      mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/volume"
	GoogleHomeMiniSmartSpeakerMQTTTopicMute        mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/mute"
	GoogleHomeMiniSmartSpeakerMQTTTopicStateVolume mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/state/volume"
	GoogleHomeMiniSmartSpeakerMQTTTopicStateMute   mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/state/mute"
)

type GoogleHomeMiniSmartSpeaker struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	mutex    sync.RWMutex
	initOnce sync.Once

	clientGoogleHome *client.GoogleHome
	clientChromecast *chromecast.Player
	host             string
}

func NewGoogleHomeMiniSmartSpeaker(host string) *GoogleHomeMiniSmartSpeaker {
	device := &GoogleHomeMiniSmartSpeaker{
		host: host,
	}
	device.Init()
	device.SetDescription("Google Home Mini")

	return device
}

func (d *GoogleHomeMiniSmartSpeaker) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeSmartSpeaker,
	}
}

func (d *GoogleHomeMiniSmartSpeaker) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 10)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Second * 30)
	taskLiveness.SetName("device-smart-home-google-home-mini-liveness")

	return []workers.Task{
		taskLiveness,
	}
}

func (d *GoogleHomeMiniSmartSpeaker) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		GoogleHomeMiniSmartSpeakerMQTTTopicVolume,
		GoogleHomeMiniSmartSpeakerMQTTTopicMute,
		GoogleHomeMiniSmartSpeakerMQTTTopicStateVolume,
		GoogleHomeMiniSmartSpeakerMQTTTopicStateMute,
	}
}

func (d *GoogleHomeMiniSmartSpeaker) ClientGoogleHome() *client.GoogleHome {
	d.mutex.RLock()
	c := d.clientGoogleHome
	d.mutex.RUnlock()

	if c != nil {
		return c
	}

	ctrl := home.NewClient(d.host)

	d.mutex.Lock()
	d.clientGoogleHome = ctrl
	d.mutex.Unlock()

	return ctrl
}

func (d *GoogleHomeMiniSmartSpeaker) ClientChromecast() *chromecast.Player {
	d.mutex.RLock()
	c := d.clientChromecast
	d.mutex.RUnlock()

	if c != nil {
		return c
	}

	ctrl := chromecast.New(d.host, chromecast.DefaultPort)

	d.mutex.Lock()
	d.clientChromecast = ctrl
	d.mutex.Unlock()

	return ctrl
}

func (d *GoogleHomeMiniSmartSpeaker) UpdateStatus(status boggart.DeviceStatus) {
	if status == boggart.DeviceStatusOffline && status != d.Status() {
		d.mutex.Lock()
		d.clientGoogleHome = nil

		if d.clientChromecast != nil {
			go d.clientChromecast.Close()
		}

		d.clientChromecast = nil
		d.mutex.Unlock()
	}

	d.DeviceBase.UpdateStatus(status)
}

func (d *GoogleHomeMiniSmartSpeaker) taskLiveness(ctx context.Context) (interface{}, error) {
	ctrl := d.ClientGoogleHome()

	response, err := ctrl.Info.GetEurekaInfo(info.NewGetEurekaInfoParams().
		WithOptions(home.EurekaInfoOptionDetail.Value()).
		WithParams(home.EurekaInfoParamDeviceInfo.Value()))
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	if d.Status() == boggart.DeviceStatusOnline {
		return nil, nil
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)

	if d.SerialNumber() == "" && response.Payload != nil {
		d.SetSerialNumber(response.Payload.DeviceInfo.MacAddress)
		d.initOnce.Do(d.initMQTTSubscribers)
	}

	return nil, nil
}

func (d *GoogleHomeMiniSmartSpeaker) initMQTTSubscribers() {
	sn := d.SerialNumberMQTTEscaped()

	d.MQTTSubscribe(GoogleHomeMiniSmartSpeakerMQTTTopicVolume.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) {
		volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
		if err != nil {
			return
		}

		d.ClientChromecast().SetVolume(volume)
	}))

	d.MQTTSubscribe(GoogleHomeMiniSmartSpeakerMQTTTopicMute.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) {
		d.ClientChromecast().SetMute(bytes.Equal(message.Payload(), []byte(`1`)))
	}))
}
