package devices

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/yryz/ds18b20"
)

const (
	DS18B20SensorMQTTTopic boggart.DeviceMQTTTopic = boggart.ComponentName + "/meter/ds18b20/+"
)

type DS18B20Sensor struct {
	lastValue int64

	boggart.DeviceBase
	boggart.DeviceSerialNumber
	boggart.DeviceMQTT

	addr string
}

func NewDS18B20Sensor(addr string) *DS18B20Sensor {
	device := &DS18B20Sensor{
		addr: addr,
	}

	device.Init()
	device.SetSerialNumber(addr)
	device.SetDescription("Sensor DS18B20 with address %s", device.SerialNumber())

	return device
}

func (d *DS18B20Sensor) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeThermometer,
	}
}

func (d *DS18B20Sensor) Ping(_ context.Context) bool {
	_, err := d.Temperature()
	return err == nil
}

func (d *DS18B20Sensor) Temperature() (float64, error) {
	return ds18b20.Temperature(d.addr)
}

func (d *DS18B20Sensor) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTask(d.taskUpdater)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(time.Minute)
	taskUpdater.SetName("device-sensor-ds18b20-updater-" + d.SerialNumber())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *DS18B20Sensor) taskUpdater(ctx context.Context) (interface{}, error) {
	if !d.IsEnabled() {
		return nil, nil
	}

	value, err := d.Temperature()
	if err != nil {
		return nil, err
	}

	prev := atomic.LoadInt64(&d.lastValue)
	current := int64(value * 1000)

	if prev != current {
		atomic.StoreInt64(&d.lastValue, current)

		d.MQTTPublishAsync(ctx, DS18B20SensorMQTTTopic.Format(d.SerialNumber()), 0, true, value)
	}

	return nil, nil
}

func (d *DS18B20Sensor) MQTTTopics() []boggart.DeviceMQTTTopic {
	return []boggart.DeviceMQTTTopic{
		DS18B20SensorMQTTTopic,
	}
}
