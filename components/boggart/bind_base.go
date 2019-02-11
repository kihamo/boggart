package boggart

import (
	"context"
	"errors"
	"sync"

	"github.com/kihamo/boggart/components/mqtt"
)

type BindBase struct {
	statusGetter BindStatusGetter
	statusSetter BindStatusSetter
	mutex        sync.RWMutex
	serialNumber string
}

func (d *BindBase) SetStatusManager(getter BindStatusGetter, setter BindStatusSetter) {
	d.mutex.Lock()
	d.statusGetter = getter
	d.statusSetter = setter
	d.mutex.Unlock()
}

func (d *BindBase) Status() BindStatus {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.statusGetter != nil {
		return d.statusGetter()
	}

	return BindStatusUnknown
}

func (d *BindBase) UpdateStatus(status BindStatus) {
	d.mutex.RLock()
	if d.statusSetter != nil {
		d.statusSetter(status)
	}
	d.mutex.RUnlock()
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

func (d *BindMQTT) MQTTPublish(ctx context.Context, topic string, payload interface{}) error {
	return d.MQTTPublishRaw(ctx, topic, 1, true, payload)
}

func (d *BindMQTT) MQTTPublishRaw(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.client == nil {
		return errors.New("MQTT client isn't init")
	}

	return d.client.Publish(ctx, topic, qos, retained, payload)
}

func (d *BindMQTT) MQTTPublishAsync(ctx context.Context, topic string, payload interface{}) error {
	return d.MQTTPublishAsyncRaw(ctx, topic, 1, true, payload)
}

func (d *BindMQTT) MQTTPublishAsyncRaw(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.client == nil {
		return errors.New("MQTT client isn't init")
	}

	d.client.PublishAsync(ctx, topic, qos, retained, payload)
	return nil
}

func CheckValueInMQTTTopic(topic string, value string, offset int) bool {
	if value == "" {
		return false
	}

	routes := mqtt.RouteSplit(topic)
	if len(routes) < offset {
		return false
	}

	return routes[len(routes)-offset] == value
}

func CheckSerialNumberInMQTTTopic(bind Bind, topic string, offset int) bool {
	return CheckValueInMQTTTopic(topic, mqtt.NameReplace(bind.SerialNumber()), offset)
}

func WrapMQTTSubscribeDeviceIsOnline(status BindStatusGetter, callback mqtt.MessageHandler) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
		if status() == BindStatusOnline {
			return callback(ctx, client, message)
		}

		return errors.New("bind isn't online")
	}
}
