package bind

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mobile"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

const (
	MegafonPhoneMQTTTopicBalance mqtt.Topic = boggart.ComponentName + "/service/megafon/+/balance"
)

type MegafonPhone struct {
	lastValue int64

	boggart.DeviceBindBase
	boggart.DeviceBindMQTT

	provider *mobile.Megafon
	interval time.Duration
}

func NewMegafonPhone(provider *mobile.Megafon, interval time.Duration) *MegafonPhone {
	device := &MegafonPhone{
		provider: provider,
		interval: interval,

		lastValue: -1,
	}
	device.Init()
	device.SetSerialNumber(device.Number())

	return device
}

func (d *MegafonPhone) Number() string {
	return d.provider.Number()
}

func (d *MegafonPhone) Balance(ctx context.Context) (float64, error) {
	return d.provider.Balance(ctx)
}

func (d *MegafonPhone) Tasks() []workers.Task {
	taskStateUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(d.interval)
	taskStateUpdater.SetName("bind-megafon-state-updater-" + d.Number())

	return []workers.Task{
		taskStateUpdater,
	}
}

func (d *MegafonPhone) taskStateUpdater(ctx context.Context) (interface{}, error) {
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

		d.MQTTPublishAsync(ctx, MegafonPhoneMQTTTopicBalance.Format(d.SerialNumber()), 0, true, value)
	}

	/*
		remainders, err := d.provider.Remainders(ctx)
		if err != nil {
			return nil, err
		}
	*/

	return nil, nil
}

func (d *MegafonPhone) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(MegafonPhoneMQTTTopicBalance.Format(d.SerialNumber())),
	}
}
