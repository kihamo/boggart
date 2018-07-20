package devices

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

var (
	metricInternetProviderSoftVideoBalance = snitch.NewGauge(boggart.ComponentName+"_device_internet_provider_softvideo_balance_rubles_total", "SoftVideo balance in rubles")
)

type SoftVideoInternet struct {
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

func (d *SoftVideoInternet) Describe(ch chan<- *snitch.Description) {
	metricInternetProviderSoftVideoBalance.With("account", d.provider.AccountID()).Describe(ch)
}

func (d *SoftVideoInternet) Collect(ch chan<- snitch.Metric) {
	metricInternetProviderSoftVideoBalance.With("account", d.provider.AccountID()).Collect(ch)
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

	metricBalance := metricInternetProviderSoftVideoBalance.With("account", d.provider.AccountID())
	balance := float64(value)

	if balance == metricBalance.Value() {
		return nil, nil
	}

	metricBalance.Set(balance)
	d.TriggerEvent(boggart.DeviceEventSoftVideoBalanceChanged, float64(value), d.provider.AccountID())

	return nil, nil
}
