package boggart

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/pborman/uuid"
)

var replacerMQTTName = strings.NewReplacer(
	":", "-",
	"/", "-",
	"_", "-",
	",", "-",
	".", "-",
)

type DeviceBase struct {
	id          string
	description atomic.Value
	status      uint64
}

func (d *DeviceBase) Init() {
	d.id = uuid.New()
	d.UpdateStatus(DeviceStatusInitializing)
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

type DeviceSerialNumber struct {
	mutex        sync.RWMutex
	serialNumber string
}

func (d *DeviceSerialNumber) SerialNumber() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.serialNumber
}

func (d *DeviceSerialNumber) SerialNumberMQTTEscaped() string {
	return mqtt.NameReplace(d.SerialNumber())
}

func (d *DeviceSerialNumber) SetSerialNumber(serialNumber string) {
	d.mutex.Lock()
	d.serialNumber = serialNumber
	d.mutex.Unlock()
}

type DeviceMQTT struct {
	Device

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

func WrapMQTTSubscribeDeviceIsOnline(device Device, callback mqtt.MessageHandler) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
		if device.Status() == DeviceStatusOnline {
			callback(ctx, client, message)
		}
	}
}
