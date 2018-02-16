package devices

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/mobile"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

var (
	metricPhoneMegafonBalance                  = snitch.NewGauge(boggart.ComponentName+"_device_phone_megafon_balance_rubles_total", "Megafon balance in rubles")
	metricPhoneMegafonUsedVoice                = snitch.NewGauge(boggart.ComponentName+"_device_phone_megafon_used_voice_minutes", "Megafon used voice in minutes")
	metricPhoneMegafonUsedSms                  = snitch.NewGauge(boggart.ComponentName+"_device_phone_megafon_used_sms", "Megafon used sms")
	metricPhoneMegafonUsedInternet             = snitch.NewGauge(boggart.ComponentName+"_device_phone_megafon_used_internet_gigabytes", "Megafon used internet in GB")
	metricPhoneMegafonUsedInternetProlongation = snitch.NewGauge(boggart.ComponentName+"_device_phone_megafon_used_internet_prolongation_gigabytes", "Megafon used internet prolongation in GB")
)

type MegafonPhone struct {
	boggart.DeviceBase

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

func (d *MegafonPhone) Describe(ch chan<- *snitch.Description) {
	number := d.Number()

	metricPhoneMegafonBalance.With("number", number).Describe(ch)
	metricPhoneMegafonUsedVoice.With("number", number).Describe(ch)
	metricPhoneMegafonUsedSms.With("number", number).Describe(ch)
	metricPhoneMegafonUsedInternet.With("number", number).Describe(ch)
	metricPhoneMegafonUsedInternetProlongation.With("number", number).Describe(ch)
}

func (d *MegafonPhone) Collect(ch chan<- snitch.Metric) {
	number := d.Number()

	metricPhoneMegafonBalance.With("number", number).Collect(ch)
	metricPhoneMegafonUsedVoice.With("number", number).Collect(ch)
	metricPhoneMegafonUsedSms.With("number", number).Collect(ch)
	metricPhoneMegafonUsedInternet.With("number", number).Collect(ch)
	metricPhoneMegafonUsedInternetProlongation.With("number", number).Collect(ch)
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

	number := d.Number()

	value, err := d.provider.Balance(ctx)
	if err != nil {
		return nil, err
	}

	metricPhoneMegafonBalance.With("number", number).Set(float64(value))

	remainders, err := d.provider.Remainders(ctx)
	if err != nil {
		return nil, err
	}

	metricPhoneMegafonUsedVoice.With("number", number).Set(remainders.Voice)
	metricPhoneMegafonUsedSms.With("number", number).Set(remainders.Sms)
	metricPhoneMegafonUsedInternet.With("number", number).Set(remainders.Internet)
	metricPhoneMegafonUsedInternetProlongation.With("number", number).Set(remainders.InternetProlongation)

	return nil, nil
}
