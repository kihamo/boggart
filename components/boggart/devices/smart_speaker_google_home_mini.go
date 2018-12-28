package devices

import (
	"context"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/google/home"
	"github.com/kihamo/boggart/components/boggart/providers/google/home/client"
	"github.com/kihamo/boggart/components/boggart/providers/google/home/client/info"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

type GoogleHomeMiniSmartSpeaker struct {
	boggart.DeviceBase
	boggart.DeviceSerialNumber

	mutex sync.RWMutex

	client *client.GoogleHome
	host   string
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

func (d *GoogleHomeMiniSmartSpeaker) Client() *client.GoogleHome {
	d.mutex.RLock()
	c := d.client
	d.mutex.RUnlock()

	if c != nil {
		return c
	}

	ctrl := home.NewClient(d.host)

	d.mutex.Lock()
	d.client = ctrl
	d.mutex.Unlock()

	return ctrl
}

func (d *GoogleHomeMiniSmartSpeaker) taskLiveness(ctx context.Context) (interface{}, error) {
	ctrl := d.Client()

	mode := "detail"
	response, err := ctrl.Info.GetEurekaInfo(info.NewGetEurekaInfoParams().WithOptions(&mode))
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	if d.Status() == boggart.DeviceStatusOnline {
		return nil, nil
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)

	if d.SerialNumber() == "" {
		d.SetSerialNumber(response.Payload.MacAddress)
	}

	return nil, nil
}
