package devices

import (
	"context"
	"errors"
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

	mutex  sync.RWMutex
	gpio   gpio.GPIOPin
	number int64
	mode   gpio.PinMode
}

func NewGPIOPin(pin int64, mode gpio.PinMode) *GPIOPin {
	device := &GPIOPin{
		number: pin,
		mode:   mode,
	}
	device.Init()
	device.SetDescription(fmt.Sprintf("GPIO pin #%d with %s mode", pin, mode))

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

func (d *GPIOPin) IsReadable() bool {
	return d.mode == gpio.PIN_IN
}

func (d *GPIOPin) IsWritable() bool {
	return d.mode == gpio.PIN_OUT
}

func (d *GPIOPin) Up() error {
	if !d.IsWritable() {
		return errors.New("GPIO isn't writable")
	}

	if pin := d.pin(); pin != nil {
		return pin.Up()
	}

	return nil
}

func (d *GPIOPin) Down() error {
	if !d.IsWritable() {
		return errors.New("GPIO isn't writable")
	}

	if pin := d.pin(); pin != nil {
		return pin.Down()
	}

	return nil
}

func (d *GPIOPin) Ping(_ context.Context) bool {
	pin := d.pin()
	if pin == nil {
		return true
	}

	_, err := pin.Mode()
	return err == nil
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

	pin, err := gpio.NewPin(d.number, d.mode)
	if err != nil {
		fmt.Printf("Init pin failed with error %s\n", err.Error())
		return nil, err, false
	}

	pin.SetCallbackChange(d.callback)

	d.mutex.Lock()
	d.gpio = pin
	d.mutex.Unlock()

	return nil, nil, true
}

func (d *GPIOPin) pin() gpio.GPIOPin {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.gpio
}

func (d *GPIOPin) callback(status bool, changedAt time.Time, prevChangedAt *time.Time) {
	if status {
		d.TriggerEvent(boggart.DeviceEventGPIOPinChanged, d.number, true, changedAt, prevChangedAt)
	} else {
		d.TriggerEvent(boggart.DeviceEventGPIOPinChanged, d.number, false, changedAt, prevChangedAt)
	}
}
