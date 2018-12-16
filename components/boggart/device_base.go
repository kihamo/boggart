package boggart

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
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

func (d *DeviceBase) SetDescription(description string, v ...interface{}) {
	d.description.Store(fmt.Sprintf(description, v...))
}

func (d *DeviceBase) Types() []DeviceType {
	return nil
}

func (d *DeviceBase) IsEnabled() bool {
	return atomic.LoadUint64(&d.enabled) == 1
}

func (d *DeviceBase) Enable() error {
	atomic.StoreUint64(&d.enabled, 1)
	d.TriggerEvent(context.TODO(), DeviceEventDeviceEnabled, d)

	return nil
}

func (d *DeviceBase) Ping(_ context.Context) bool {
	return false
}

func (d *DeviceBase) Disable() error {
	atomic.StoreUint64(&d.enabled, 0)
	d.TriggerEvent(context.TODO(), DeviceEventDeviceDisabled, d)

	return nil
}

func (d *DeviceBase) TriggerEventChannel() <-chan DeviceTriggerEvent {
	return d.triggerEventsChannel
}

func (d *DeviceBase) TriggerEvent(ctx context.Context, event workers.Event, args ...interface{}) {
	if d.triggerEventsChannel == nil {
		return
	}

	go func() {
		d.triggerEventsChannel <- NewDeviceTriggerEventBase(ctx, event, append([]interface{}{d}, args...))
	}()
}

type DeviceSerialNumber struct {
	mutex        sync.RWMutex
	serialNumber string
}

func (d *DeviceSerialNumber) SerialNumber() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.serialNumber
}

func (d *DeviceSerialNumber) SetSerialNumber(serialNumber string) {
	d.mutex.Lock()
	d.serialNumber = serialNumber
	d.mutex.Unlock()
}

type DeviceMQTT struct {
	mutex  sync.RWMutex
	client mqtt.Component
}

func (d *DeviceMQTT) SetMQTTClient(client mqtt.Component) {
	d.mutex.Lock()
	d.client = client
	d.mutex.Unlock()
}

func (d *DeviceMQTT) MQTTPublish(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.client == nil {
		return errors.New("MQTT client isn't init")
	}

	return d.client.Publish(ctx, topic, qos, retained, payload)
}
