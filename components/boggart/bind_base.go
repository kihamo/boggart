package boggart

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
)

type BindBase struct {
	mutex sync.RWMutex

	serialNumber string
	status       uint64
}

func (d *BindBase) Init() {
	d.UpdateStatus(BindStatusInitializing)
}

func (d *BindBase) Status() BindStatus {
	return BindStatus(atomic.LoadUint64(&d.status))
}

func (d *BindBase) UpdateStatus(status BindStatus) {
	atomic.StoreUint64(&d.status, uint64(status))
}

func (d *BindBase) SerialNumber() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.serialNumber
}

func (d *BindBase) SetSerialNumber(serialNumber string) {
	d.mutex.Lock()
	d.serialNumber = serialNumber
	d.mutex.Unlock()
}

type BindMQTT struct {
	mutex  sync.RWMutex
	client mqtt.Component
}

func (d *BindMQTT) SetMQTTClient(client mqtt.Component) {
	d.mutex.Lock()
	d.client = client
	d.mutex.Unlock()
}

func (d *BindMQTT) MQTTPublish(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.client == nil {
		return errors.New("MQTT client isn't init")
	}

	return d.client.Publish(ctx, topic, qos, retained, payload)
}

func (d *BindMQTT) MQTTPublishAsync(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.client == nil {
		return errors.New("MQTT client isn't init")
	}

	d.client.PublishAsync(ctx, topic, qos, retained, payload)
	return nil
}

func CheckSerialNumberInMQTTTopic(bind Bind, topic string, offset int) bool {
	sn := mqtt.NameReplace(bind.SerialNumber())

	if sn == "" {
		return false
	}

	routes := mqtt.RouteSplit(topic)
	if len(routes) < offset {
		return false
	}

	return routes[len(routes)-offset] == sn
}

func WrapMQTTSubscribeDeviceIsOnline(bind Bind, callback mqtt.MessageHandler) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
		if bind.Status() == BindStatusOnline {
			return callback(ctx, client, message)
		}

		return errors.New("bind isn't online")
	}
}
