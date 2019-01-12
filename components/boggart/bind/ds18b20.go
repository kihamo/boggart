package bind

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/yryz/ds18b20"
)

const (
	DS18B20MQTTTopic mqtt.Topic = boggart.ComponentName + "/meter/ds18b20/+"
)

type DS18B20 struct {
	lastValue int64

	boggart.DeviceBindBase
	boggart.DeviceBindMQTT
}

type DS18B20Config struct {
	Address string `valid:"required"`
}

func (d DS18B20) Config() interface{} {
	return &DS18B20Config{}
}

func (d DS18B20) CreateBind(c interface{}) (boggart.DeviceBind, error) {
	config := c.(*DS18B20Config)

	device := &DS18B20{
		lastValue: -1,
	}
	device.Init()
	device.SetSerialNumber(config.Address)

	return device, nil
}

func (d *DS18B20) Tasks() []workers.Task {
	sn := d.SerialNumber()

	taskLiveness := task.NewFunctionTask(d.taskLiveness)
	taskLiveness.SetTimeout(time.Second * 5)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(time.Minute)
	taskLiveness.SetName("bind-ds18b20-liveness-" + sn)

	taskStateUpdater := task.NewFunctionTask(d.taskStateUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(time.Minute)
	taskStateUpdater.SetName("bind-ds18b20-state-updater-" + sn)

	return []workers.Task{
		taskLiveness,
		taskStateUpdater,
	}
}

func (d *DS18B20) taskLiveness(ctx context.Context) (interface{}, error) {
	devices, err := ds18b20.Sensors()
	if err != nil {
		d.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	sn := d.SerialNumber()

	for _, device := range devices {
		if device == sn {
			d.UpdateStatus(boggart.DeviceStatusOnline)
			return nil, nil
		}
	}

	d.UpdateStatus(boggart.DeviceStatusOffline)
	return nil, nil
}

func (d *DS18B20) taskStateUpdater(ctx context.Context) (interface{}, error) {
	if d.Status() != boggart.DeviceStatusOnline {
		return nil, nil
	}

	sn := d.SerialNumber()

	value, err := ds18b20.Temperature(sn)
	if err != nil {
		return nil, err
	}

	prev := atomic.LoadInt64(&d.lastValue)
	current := int64(value * 1000)

	if prev != current {
		atomic.StoreInt64(&d.lastValue, current)

		d.MQTTPublishAsync(ctx, DS18B20MQTTTopic.Format(sn), 0, true, value)
	}

	return nil, nil
}

func (d *DS18B20) MQTTTopics() []mqtt.Topic {
	return []mqtt.Topic{
		mqtt.Topic(DS18B20MQTTTopic.Format(mqtt.NameReplace(d.SerialNumber()))),
	}
}
