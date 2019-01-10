package boggart

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
)

type DeviceBindBase struct {
	status uint64
}

func (d *DeviceBindBase) Init() {
	d.UpdateStatus(DeviceStatusInitializing)
}

func (d *DeviceBindBase) Status() DeviceStatus {
	return DeviceStatus(atomic.LoadUint64(&d.status))
}

func (d *DeviceBindBase) UpdateStatus(status DeviceStatus) {
	atomic.StoreUint64(&d.status, uint64(status))
}

type DeviceBindSerialNumber struct {
	mutex        sync.RWMutex
	serialNumber string
}

func (d *DeviceBindSerialNumber) SerialNumber() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.serialNumber
}

func (d *DeviceBindSerialNumber) SerialNumberMQTTEscaped() string {
	return mqtt.NameReplace(d.SerialNumber())
}

func (d *DeviceBindSerialNumber) SetSerialNumber(serialNumber string) {
	d.mutex.Lock()
	d.serialNumber = serialNumber
	d.mutex.Unlock()
}

type DeviceBindMQTT struct {
	Device

	mutex  sync.RWMutex
	client mqtt.Component
}

func (d *DeviceBindMQTT) SetMQTTClient(client mqtt.Component) {
	d.mutex.Lock()
	d.client = client
	d.mutex.Unlock()
}

func (d *DeviceBindMQTT) MQTTPublish(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error {
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

func (d *DeviceBindMQTT) MQTTPublishAsync(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) {
	go func() {
		d.MQTTPublish(ctx, topic, qos, retained, payload)
	}()
}

func (d *DeviceBindMQTT) MQTTSubscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.client == nil {
		return errors.New("MQTT client isn't init")
	}

	return d.client.Subscribe(topic, qos, callback)
}

func WrapMQTTSubscribeDeviceIsOnline(bind DeviceBind, callback mqtt.MessageHandler) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
		if bind.Status() == DeviceStatusOnline {
			callback(ctx, client, message)
		}
	}
}
