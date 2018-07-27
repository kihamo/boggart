package devices

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/gpio"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

type DoorGPIOPin struct {
	boggart.DeviceBase

	mutex     sync.RWMutex
	pinGPIO   gpio.GPIOPin
	pinNumber int64
}

func NewGPIOPin(pin int64) *DoorGPIOPin {
	device := &DoorGPIOPin{
		pinNumber: pin,
	}
	device.Init()
	device.SetDescription(fmt.Sprintf("GPIO pin #%d", pin))

	return device
}

func (d *DoorGPIOPin) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeGPIOPin,
	}
}

func (d *DoorGPIOPin) IsOn() bool {
	return !d.IsOff()
}

func (d *DoorGPIOPin) IsOff() bool {
	if pin := d.pin(); pin != nil {
		return pin.Status()
	}

	return false
}

func (d *DoorGPIOPin) Ping(_ context.Context) bool {
	return true
	//return d.pin() != nil
}

func (d *DoorGPIOPin) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTillStopTask(d.taskPin)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(time.Second)
	taskUpdater.SetName("device-gpio-pin-" + d.Id())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *DoorGPIOPin) taskPin(ctx context.Context) (interface{}, error, bool) {
	if !d.IsEnabled() {
		return nil, nil, false
	}

	pin, err := gpio.NewPin(d.pinNumber, gpio.PIN_IN)
	if err != nil {
		return nil, err, false
	}

	pin.SetCallbackChange(d.callback)

	d.mutex.Lock()
	d.pinGPIO = pin
	d.mutex.Unlock()

	return nil, nil, true
}

func (d *DoorGPIOPin) pin() gpio.GPIOPin {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.pinGPIO
}

func (d *DoorGPIOPin) callback(status bool, changed *time.Time) {
	if status {
		d.TriggerEvent(boggart.DeviceEventGPIOPinChanged, d.pinNumber, true, changed)
	} else {
		d.TriggerEvent(boggart.DeviceEventGPIOPinChanged, d.pinNumber, false, changed)
	}
}
