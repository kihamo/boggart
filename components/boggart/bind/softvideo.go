package bind

import (
	"context"
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
	boggart.DeviceBindMQTT

	provider *softvideo.Client
}

type SoftVideoConfig struct {
	Login    string `valid:"required"`
	Password string `valid:"required"`
}

func (d SoftVideo) Config() interface{} {
	return &SoftVideoConfig{}
}

func (d SoftVideo) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*SoftVideoConfig)

	device := &SoftVideo{
		provider:  softvideo.NewClient(config.Login, config.Password),
		lastValue: -1,
	}
	device.Init()
	device.SetSerialNumber(config.Login)

	return device, nil
}

func (d *SoftVideo) Balance(ctx context.Context) (float64, error) {
	return d.provider.Balance(ctx)
}

func (d *SoftVideo) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(time.Hour)
	taskUpdater.SetName("bind-softvideo-updater-" + d.provider.AccountID())

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
