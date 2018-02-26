package devices

import (
	"context"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/apcupsd"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

type ApcupsdUPS struct {
	boggart.DeviceWithSerialNumber

	client   *apcupsd.Client
	interval time.Duration
}

func NewApcupsdUPS(client *apcupsd.Client, interval time.Duration) *ApcupsdUPS {
	device := &ApcupsdUPS{
		client:   client,
		interval: interval,
	}
	device.Init()
	device.SetDescription("UPS")

	return device
}

func (d *ApcupsdUPS) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeUPS,
	}
}

func (d *ApcupsdUPS) Describe(ch chan<- *snitch.Description) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}
}

func (d *ApcupsdUPS) Collect(ch chan<- snitch.Metric) {
	serialNumber := d.SerialNumber()
	if serialNumber == "" {
		return
	}
}

func (d *ApcupsdUPS) Ping(ctx context.Context) bool {
	status, err := d.client.Status(ctx)
	if err == nil && status.Status != nil {
		return *status.Status == "ONLINE"
	}

	return false
}

func (d *ApcupsdUPS) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillStopTask(d.taskSerialNumber)
	taskSerialNumber.SetTimeout(time.Second * 5)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Minute)
	taskSerialNumber.SetName("device-ups-apcupsd-serial-number")

	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(d.interval)
	taskUpdater.SetName("device-ups-apcupsd-updater-" + d.Id())

	return []workers.Task{
		taskSerialNumber,
		taskUpdater,
	}
}

func (d *ApcupsdUPS) taskSerialNumber(ctx context.Context) (interface{}, error, bool) {
	if !d.IsEnabled() {
		return nil, nil, false
	}

	status, err := d.client.Status(ctx)
	if err != nil || status.SerialNumber == nil {
		return nil, err, false
	}

	d.SetSerialNumber(*status.SerialNumber)
	d.SetDescription("UPS with serial number " + *status.SerialNumber)

	return nil, nil, true
}

func (d *ApcupsdUPS) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	// TODO: автоматически подстроится под интервал обновления на apcupsd, что бы лишний раз не гонять тикет

	return nil, nil
}
