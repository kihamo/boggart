package devices

import (
	"context"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/gpio"
	"github.com/kihamo/go-workers/event"
	"github.com/kihamo/snitch"
)

var (
	EventDoorGPIOReedSwitchOpen  = event.NewBaseEvent(boggart.ComponentName + "_device_door_gpio_reed_switch_open")
	EventDoorGPIOReedSwitchClose = event.NewBaseEvent(boggart.ComponentName + "_device_door_gpio_reed_switch_close")

	metricDoorGPIOReedSwitchStatus = snitch.NewGauge(boggart.ComponentName+"_device_door_gpio_reed_switch_status", "Door status")
)

type DoorGPIOReedSwitch struct {
	boggart.DeviceBase

	pin gpio.GPIOPin
}

func NewDoorGPIOReedSwitch(pin int64) (*DoorGPIOReedSwitch, error) {
	p, err := gpio.NewPin(pin, gpio.PIN_IN)
	if err != nil {
		return nil, err
	}

	device := &DoorGPIOReedSwitch{
		pin: p,
	}
	device.Init()
	device.SetDescription("Door GPIO reed switch")

	p.SetCallbackChange(device.callback)

	return device, nil
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
	return d.pin.Status()
}

func (d *DoorGPIOReedSwitch) Describe(ch chan<- *snitch.Description) {
	metricDoorGPIOReedSwitchStatus.With("pin", strconv.FormatInt(d.pin.Number(), 10)).Describe(ch)
}

func (d *DoorGPIOReedSwitch) Collect(ch chan<- snitch.Metric) {
	metricStatus := metricDoorGPIOReedSwitchStatus.With("pin", strconv.FormatInt(d.pin.Number(), 10))

	if d.IsOpen() {
		metricStatus.Set(1)
	} else {
		metricStatus.Set(0)
	}

	metricStatus.Collect(ch)
}

func (d *DoorGPIOReedSwitch) Ping(_ context.Context) bool {
	return true
}

func (d *DoorGPIOReedSwitch) callback(status bool, changed *time.Time) {
	if status {
		d.TriggerEvent(EventDoorGPIOReedSwitchClose, changed)
	} else {
		d.TriggerEvent(EventDoorGPIOReedSwitchOpen, changed)
	}
}
