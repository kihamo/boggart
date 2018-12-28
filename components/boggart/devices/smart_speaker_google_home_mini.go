package devices

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/barnybug/go-cast"
	"github.com/barnybug/go-cast/controllers"
	"github.com/barnybug/go-cast/events"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/google/home"
	"github.com/kihamo/boggart/components/boggart/providers/google/home/client"
	"github.com/kihamo/boggart/components/boggart/providers/google/home/client/info"
	"github.com/kihamo/boggart/components/mqtt"
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
	clientChromecast *cast.Client
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

func (d *GoogleHomeMiniSmartSpeaker) ClientChromecast() (*cast.Client, error) {
	d.mutex.RLock()
	c := d.clientChromecast
	d.mutex.RUnlock()

	if c != nil {
		return c, nil
	}

	ctrl := cast.NewClient(net.ParseIP(d.host), 8009)
	err := ctrl.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	// TODO: fix quit
	go func() {
		sn := d.SerialNumberMQTTEscaped()
		ctx := context.Background()

		for event := range ctrl.Events {
			switch t := event.(type) {
			case events.StatusUpdated:
				d.MQTTPublishAsync(ctx, GoogleHomeMiniSmartSpeakerMQTTTopicStateVolume.Format(sn), 2, true, int64(math.Round(t.Level*100)))
				d.MQTTPublishAsync(ctx, GoogleHomeMiniSmartSpeakerMQTTTopicStateMute.Format(sn), 2, true, t.Muted)

			default:
				fmt.Println("FIXME GHM:", t)
			}
		}
	}()

	d.mutex.Lock()
	d.clientChromecast = ctrl
	d.mutex.Unlock()

	return ctrl, nil
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
		volume, err := strconv.ParseFloat(string(message.Payload()), 64)
		if err != nil {
			return
		}

		if volume > 100 {
			volume = 100
		} else if volume < 0 {
			volume = 0
		}

		volume /= 100

		ctrl, err := d.ClientChromecast()
		if err != nil {
			return
		}

		receiver := ctrl.Receiver()

		// иначе делает двойной запрос на установку к колонке
		response, err := receiver.GetVolume(ctx)
		if err != nil {
			return
		}

		_, _ = receiver.SetVolume(ctx, &controllers.Volume{
			Level: &volume,
			Muted: response.Muted,
		})
	}))

	d.MQTTSubscribe(GoogleHomeMiniSmartSpeakerMQTTTopicMute.Format(sn), 0, boggart.WrapMQTTSubscribeDeviceIsOnline(d, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) {
		mute := bytes.Equal(message.Payload(), []byte(`1`))

		ctrl, err := d.ClientChromecast()
		if err != nil {
			return
		}

		receiver := ctrl.Receiver()
		_, _ = receiver.SetVolume(ctx, &controllers.Volume{
			Muted: &mute,
		})
	}))
}
