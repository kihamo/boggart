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

type GPIOPin struct {
	boggart.DeviceBase

	mutex     sync.RWMutex
	pinGPIO   gpio.GPIOPin
	pinNumber int64
}

func NewGPIOPin(pin int64) *GPIOPin {
	device := &GPIOPin{
		pinNumber: pin,
	}
	device.Init()
	device.SetDescription(fmt.Sprintf("GPIO pin #%d", pin))

	return device
}

func (d *GPIOPin) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeGPIOPin,
	}
}

func (d *GPIOPin) IsOn() bool {
	return !d.IsOff()
}

func (d *GPIOPin) IsOff() bool {
	if pin := d.pin(); pin != nil {
		return pin.Status()
	}

	return false
}

func (d *GPIOPin) Ping(_ context.Context) bool {
	return true
	//return d.pin() != nil
}

func (d *GPIOPin) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTillStopTask(d.taskPin)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(time.Second)
	taskUpdater.SetName("device-gpio-pin-" + d.Id())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *GPIOPin) taskPin(ctx context.Context) (interface{}, error, bool) {
	if !d.IsEnabled() {
		return nil, nil, false
	}

	pin, err := gpio.NewPin(d.pinNumber, gpio.PIN_IN)
	if err != nil {
		fmt.Printf("Init pin failed with error %s\n", err.Error())
		return nil, err, false
	}

	pin.SetCallbackChange(d.callback)

	d.mutex.Lock()
	d.pinGPIO = pin
	d.mutex.Unlock()

	return nil, nil, true
}

func (d *GPIOPin) pin() gpio.GPIOPin {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.pinGPIO
}

func (d *GPIOPin) callback(status bool, changedAt time.Time, prevChangedAt *time.Time) {
	if status {
		d.TriggerEvent(boggart.DeviceEventGPIOPinChanged, d.pinNumber, true, changedAt, prevChangedAt)
	} else {
		d.TriggerEvent(boggart.DeviceEventGPIOPinChanged, d.pinNumber, false, changedAt, prevChangedAt)
	}
}
