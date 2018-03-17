package devices

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/gpio"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/event"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/snitch"
)

var (
	EventDoorGPIOReedSwitchOpen  = event.NewBaseEvent("DeviceDoorGPIOReedSwitchOpen")
	EventDoorGPIOReedSwitchClose = event.NewBaseEvent("DeviceDoorGPIOReedSwitchClose")

	metricDoorGPIOReedSwitchStatus = snitch.NewGauge(boggart.ComponentName+"_device_door_gpio_reed_switch_status", "Door status")
)

type DoorGPIOReedSwitch struct {
	boggart.DeviceBase

	mutex     sync.RWMutex
	pinGPIO   gpio.GPIOPin
	pinNumber int64
}

func NewDoorGPIOReedSwitch(pin int64) *DoorGPIOReedSwitch {
	device := &DoorGPIOReedSwitch{
		pinNumber: pin,
	}
	device.Init()
	device.SetDescription("Door GPIO reed switch")

	return device
}

func (d *DoorGPIOReedSwitch) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeDoor,
	}
}

func (d *DoorGPIOReedSwitch) IsOpen() bool {
	return !d.IsClose()
}

func (d *DoorGPIOReedSwitch) IsClose() bool {
	pin := d.pin()
	if pin != nil {
		return pin.Status()
	}

	return false
}

func (d *DoorGPIOReedSwitch) Describe(ch chan<- *snitch.Description) {
	pin := d.pin()
	if pin == nil {
		return
	}

	metricDoorGPIOReedSwitchStatus.With("pin", strconv.FormatInt(pin.Number(), 10)).Describe(ch)
}

func (d *DoorGPIOReedSwitch) Collect(ch chan<- snitch.Metric) {
	pin := d.pin()
	if pin == nil {
		return
	}

	metricStatus := metricDoorGPIOReedSwitchStatus.With("pin", strconv.FormatInt(pin.Number(), 10))

	if d.IsOpen() {
		metricStatus.Set(1)
	} else {
		metricStatus.Set(0)
	}

	metricStatus.Collect(ch)
}

func (d *DoorGPIOReedSwitch) Ping(_ context.Context) bool {
	return d.pin() != nil
}

func (d *DoorGPIOReedSwitch) Tasks() []workers.Task {
	taskUpdater := task.NewFunctionTillStopTask(d.taskPin)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(time.Second)
	taskUpdater.SetName("device-door-gpio-reed-switch-pin-" + d.Id())

	return []workers.Task{
		taskUpdater,
	}
}

func (d *DoorGPIOReedSwitch) taskPin(ctx context.Context) (interface{}, error, bool) {
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

	d.SetDescription(d.Description() + " with PIN #" + strconv.FormatInt(pin.Number(), 10))

	return nil, nil, true
}

func (d *DoorGPIOReedSwitch) pin() gpio.GPIOPin {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.pinGPIO
}

func (d *DoorGPIOReedSwitch) callback(status bool, changed *time.Time) {
	if status {
		d.TriggerEvent(EventDoorGPIOReedSwitchClose, changed)
	} else {
		d.TriggerEvent(EventDoorGPIOReedSwitchOpen, changed)
	}
}
