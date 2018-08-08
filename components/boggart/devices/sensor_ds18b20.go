package devices

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"periph.io/x/periph/conn/onewire"
	"periph.io/x/periph/devices/ds18b20"
)

type DS18B20Sensor struct {
	lastValue int64

	boggart.DeviceBase
	boggart.DeviceSerialNumber

	device *ds18b20.Dev
}

func NewDS18B20Sensor(o onewire.Bus, addr onewire.Address, resolutionBits int) *DS18B20Sensor {
	dev, _ := ds18b20.New(o, addr, resolutionBits)

	device := &DS18B20Sensor{
		device: dev,
	}

	device.Init()
	device.SetSerialNumber(fmt.Sprintf("0x%02x", uint64(addr)))
	device.SetDescription("Sensor DS18B20 with address %s", device.SerialNumber())

	return device
}

func (d *DS18B20Sensor) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeThermometer,
	}
}

func (d *DS18B20Sensor) Ping(_ context.Context) bool {
	_, err := d.device.LastTemp()
	return err == nil
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

	value, err := d.device.LastTemp()
	if err != nil {
		return nil, err
	}

	prev := atomic.LoadInt64(&d.lastValue)
	current := int64(value)

	if prev != current {
		atomic.StoreInt64(&d.lastValue, current)
		d.TriggerEvent(boggart.DeviceEventDS18B20Changed, current, d.SerialNumber())
	}

	return nil, nil
}
