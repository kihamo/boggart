package devices

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
	SoftVideoInternetMQTTTopicBalance mqtt.Topic = boggart.ComponentName + "/service/softvideo/+/balance"
)

type SoftVideoInternet struct {
	lastValue int64

	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	provider *softvideo.Client
	interval time.Duration
}

func NewSoftVideoInternet(provider *softvideo.Client, interval time.Duration) *SoftVideoInternet {
	device := &SoftVideoInternet{
		provider: provider,
		interval: interval,

		lastValue: -1,
	}
	device.Init()
	device.SetSerialNumber(provider.AccountID())
	device.SetDescription("SoftVideo internet provider")

	return device
}

func (d *SoftVideoInternet) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeInternetProvider,
	}
}

func (d *SoftVideoInternet) Balance(ctx context.Context) (float64, error) {
	return d.provider.Balance(ctx)
}

func (d *SoftVideoInternet) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-internet-provider-softvideo-updater-" + d.provider.AccountID())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *SoftVideoInternet) taskUpdater(ctx context.Context) (interface{}, error) {
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

		d.MQTTPublishAsync(ctx, SoftVideoInternetMQTTTopicBalance.Format(d.SerialNumber()), 0, true, value)
	}

	return nil, nil
}

func (d *SoftVideoInternet) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(SoftVideoInternetMQTTTopicBalance.Format(d.SerialNumber())),
	}
}
