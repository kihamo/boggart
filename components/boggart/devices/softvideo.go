package devices

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

const (
	SoftVideoMQTTTopicBalance mqtt.Topic = boggart.ComponentName + "/service/softvideo/+/balance"
)

type SoftVideo struct {
	lastValue int64

	boggart.DeviceBindBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	provider *softvideo.Client
}

func (d SoftVideo) CreateBind(config map[string]interface{}) (boggart.DeviceBind, error) {
	login, ok := config["login"]
	if !ok {
		return nil, errors.New("config option login isn't set")
	}

	if login == "" {
		return nil, errors.New("config option login is empty")
	}

	password, ok := config["password"]
	if !ok {
		return nil, errors.New("config option password isn't set")
	}

	if password == "" {
		return nil, errors.New("config option password is empty")
	}

	device := &SoftVideo{
		provider:  softvideo.NewClient(login.(string), password.(string)),
		lastValue: -1,
	}
	device.Init()
	device.SetSerialNumber(login.(string))

	return device, nil
}

func (d *SoftVideo) Balance(ctx context.Context) (float64, error) {
	return d.provider.Balance(ctx)
}

func (d *SoftVideo) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(time.Hour)
	taskUpdater.SetName("device-internet-provider-softvideo-updater-" + d.provider.AccountID())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *SoftVideo) taskUpdater(ctx context.Context) (interface{}, error) {
	value, err := d.provider.Balance(ctx)
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	d.UpdateStatus(boggart.DeviceStatusOnline)

	current := int64(value * 100)
	prev := atomic.LoadInt64(&d.lastValue)

	if current != prev {
		atomic.StoreInt64(&d.lastValue, current)

		d.MQTTPublishAsync(ctx, SoftVideoMQTTTopicBalance.Format(d.SerialNumber()), 0, true, value)
	}

	return nil, nil
}

func (d *SoftVideo) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(SoftVideoMQTTTopicBalance.Format(d.SerialNumber())),
	}
}
