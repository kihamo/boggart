package devices

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

type SoftVideoInternet struct {
	lastValue int64

	boggart.DeviceBase

	provider *softvideo.Client
	interval time.Duration
}

func NewSoftVideoInternet(provider *softvideo.Client, interval time.Duration) *SoftVideoInternet {
	device := &SoftVideoInternet{
		provider: provider,
		interval: interval,
	}
	device.Init()
	device.SetDescription("SoftVideo internet provider for account " + provider.AccountID())

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

func (d *SoftVideoInternet) Ping(_ context.Context) bool {
	return true
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
		d.TriggerEvent(boggart.DeviceEventSoftVideoBalanceChanged, value, d.provider.AccountID())
	}

	return nil, nil
}
