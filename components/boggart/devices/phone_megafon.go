package devices

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

	boggart.DeviceBase
	boggart.DeviceMQTT

	provider *mobile.Megafon
	interval time.Duration
}

func NewMegafonPhone(provider *mobile.Megafon, interval time.Duration) *MegafonPhone {
	device := &MegafonPhone{
		provider: provider,
		interval: interval,
	}
	device.Init()
	device.SetDescription("Mobile phone " + device.Number() + " of Megafon")

	return device
}

func (d *MegafonPhone) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypePhone,
		boggart.DeviceTypeInternetProvider,
	}
}

func (d *MegafonPhone) Number() string {
	return d.provider.Number()
}

func (d *MegafonPhone) Balance(ctx context.Context) (float64, error) {
	return d.provider.Balance(ctx)
}

func (d *MegafonPhone) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-phone-megafon-updater-" + d.Number())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *MegafonPhone) Ping(_ context.Context) bool {
	return true
}

func (d *MegafonPhone) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	value, err := d.provider.Balance(ctx)
	if err != nil {
		return nil, err
	}

	current := int64(value * 100)
	prev := atomic.LoadInt64(&d.lastValue)

	if current != prev {
		atomic.StoreInt64(&d.lastValue, current)

		d.MQTTPublishAsync(ctx, MegafonPhoneMQTTTopicBalance.Format(d.Number()), 0, true, value)
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
		MegafonPhoneMQTTTopicBalance,
	}
}
