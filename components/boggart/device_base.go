package boggart

import (
	"context"
	"sync/atomic"

	"github.com/kihamo/go-workers"
	"github.com/pborman/uuid"
)

type DeviceBase struct {
	id                   atomic.Value
	description          atomic.Value
	enabled              uint64
	triggerEventsChannel chan DeviceTriggerEvent
}

func (d *DeviceBase) Init() {
	d.triggerEventsChannel = make(chan DeviceTriggerEvent)
	d.SetId(uuid.New())
	d.Enable()
}

func (d *DeviceBase) Id() string {
	var id string

	if value := d.id.Load(); value != nil {
		id = value.(string)
	}

	return id
}

func (d *DeviceBase) SetId(id string) {
	d.id.Store(id)
}

func (d *DeviceBase) Description() string {
	var description string

	if value := d.description.Load(); value != nil {
		description = value.(string)
	}

	return description
}

func (d *DeviceBase) SetDescription(description string) {
	d.description.Store(description)
}

func (d *DeviceBase) Types() []DeviceType {
	return nil
}

func (d *DeviceBase) IsEnabled() bool {
	return atomic.LoadUint64(&d.enabled) == 1
}

func (d *DeviceBase) Enable() error {
	atomic.StoreUint64(&d.enabled, 1)
	d.TriggerEvent(DeviceEventDeviceEnabled, d)

	return nil
}

func (d *DeviceBase) Ping(_ context.Context) bool {
	return false
}

func (d *DeviceBase) Disable() error {
	atomic.StoreUint64(&d.enabled, 0)
	d.TriggerEvent(DeviceEventDeviceDisabled, d)

	return nil
}

func (d *DeviceBase) Listeners() []workers.ListenerWithEvents {
	return nil
}

func (d *DeviceBase) Tasks() []workers.Task {
	return nil
}

func (d *DeviceBase) TriggerEventChannel() <-chan DeviceTriggerEvent {
	return d.triggerEventsChannel
}

func (d *DeviceBase) TriggerEvent(event workers.Event, args ...interface{}) {
	if d.triggerEventsChannel == nil {
		return
	}

	go func() {
		d.triggerEventsChannel <- NewDeviceTriggerEventBase(event, append([]interface{}{d}, args...))
	}()
}
