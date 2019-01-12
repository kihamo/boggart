package bind

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
	GoogleHomeMiniMQTTTopicVolume      mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/volume"
	GoogleHomeMiniMQTTTopicMute        mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/mute"
	GoogleHomeMiniMQTTTopicStateVolume mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/state/volume"
	GoogleHomeMiniMQTTTopicStateMute   mqtt.Topic = boggart.ComponentName + "/smart-speaker/+/state/mute"
)

type GoogleHomeMini struct {
	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	mutex    sync.RWMutex
	initOnce sync.Once

	clientGoogleHome *client.GoogleHome
	clientChromecast *chromecast.Player
	host             string
}

type GoogleHomeMiniConfig struct {
	Host string `valid:"host,required"`
}

func (d GoogleHomeMini) Config() interface{} {
	return &GoogleHomeMiniConfig{}
}

func (d GoogleHomeMini) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*GoogleHomeMiniConfig)

	device := &GoogleHomeMini{
		host: config.Host,
	}
	device.Init()

	return device, nil
}

func (d *GoogleHomeMini) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 10)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Second * 30)
	taskLiveness.SetName("bind-google-home-mini-liveness")

	return []workers.Task{
		taskLiveness,
	}
}

func (d *GoogleHomeMini) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		GoogleHomeMiniMQTTTopicVolume,
		GoogleHomeMiniMQTTTopicMute,
		GoogleHomeMiniMQTTTopicStateVolume,
		GoogleHomeMiniMQTTTopicStateMute,
	}
}

func (d *GoogleHomeMini) ClientGoogleHome() *client.GoogleHome {
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

func (d *GoogleHomeMini) ClientChromecast() *chromecast.Player {
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

func (d *GoogleHomeMini) UpdateStatus(status boggart.DeviceStatus) {
	if status == boggart.DeviceStatusOffline && status != d.Status() {
		d.mutex.Lock()
		d.clientGoogleHome = nil

		if d.clientChromecast != nil {
			go d.clientChromecast.Close()
		}

		d.clientChromecast = nil
		d.mutex.Unlock()
	}

	d.DeviceBindBase.UpdateStatus(status)
}

func (d *GoogleHomeMini) taskLiveness(ctx context.Context) (interface{}, error) {
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

	if d.SerialNumber() == "" && response.Payload != nil {
		d.SetSerialNumber(response.Payload.DeviceInfo.MacAddress)
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)
	return nil, nil
}

func (d *GoogleHomeMini) MQTTSubscribers() []mqtt.Subscriber {
	return []mqtt.Subscriber{
		mqtt.NewSubscriber(GoogleHomeMiniMQTTTopicVolume.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) {
			if !d.CheckSerialNumberInMQTTTopic(message.Topic(), 2) {
				return
			}

			volume, err := strconv.ParseInt(string(message.Payload()), 10, 64)
			if err != nil {
				return
			}

			d.ClientChromecast().SetVolume(volume)
		})),
		mqtt.NewSubscriber(GoogleHomeMiniMQTTTopicMute.String(), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) {
			if !d.CheckSerialNumberInMQTTTopic(message.Topic(), 2) {
				return
			}

			d.ClientChromecast().SetMute(bytes.Equal(message.Payload(), []byte(`1`)))
		})),
	}
}
