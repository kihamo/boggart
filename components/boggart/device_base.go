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
	id                   string
	description          atomic.Value
	status               uint64
	enabled              uint64
	triggerEventsChannel chan DeviceTriggerEvent
}

func (d *DeviceBase) Init() {
	d.triggerEventsChannel = make(chan DeviceTriggerEvent)
	d.id = uuid.New()
	d.UpdateStatus(DeviceStatusInitializing)
	d.Enable()
}

func (d *DeviceBase) Id() string {
	return d.id
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

func (d *DeviceBase) Status() DeviceStatus {
	return DeviceStatus(atomic.LoadUint64(&d.status))
}

func (d *DeviceBase) UpdateStatus(status DeviceStatus) {
	atomic.StoreUint64(&d.status, uint64(status))
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

	switch value := payload.(type) {
	case string, []byte:
		// skip
	case float64, float32:
		payload = fmt.Sprintf("%.2f", value)
	case uint64, uint32, uint16, uint8, uint, int64, int32, int16, int8, int:
		payload = fmt.Sprintf("%d", value)
	case bool:
		if value {
			payload = []byte(`1`)
		} else {
			payload = []byte(`0`)
		}
	default:
		payload = fmt.Sprintf("%s", payload)
	}

	return d.client.Publish(ctx, topic, qos, retained, payload)
}

func (d *DeviceMQTT) MQTTPublishAsync(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) {
	go func() {
		d.MQTTPublish(ctx, topic, qos, retained, payload)
	}()
}

func (d *DeviceMQTT) MQTTSubscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.client == nil {
		return errors.New("MQTT client isn't init")
	}

	return d.client.Subscribe(topic, qos, callback)
}
