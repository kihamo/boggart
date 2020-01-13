package boggart

import (
	"context"
	"errors"
	"sync"

	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/shadow/components/logging"
)

type BindBase struct {
	statusManager BindStatusManager
	mutex         sync.RWMutex
	logger        logging.Logger
	serialNumber  string
}

func (b *BindBase) Run() error {
	return nil
}

func (b *BindBase) SetStatusManager(manager BindStatusManager) {
	b.mutex.Lock()
	b.statusManager = manager
	b.mutex.Unlock()
}

func (b *BindBase) Logger() logging.Logger {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.logger == nil {
		b.logger = logging.DefaultLogger()
	}

	return b.logger
}

func (b *BindBase) SetLogger(logger logging.Logger) {
	b.mutex.Lock()
	b.logger = logger
	b.mutex.Unlock()
}

func (b *BindBase) Status() BindStatus {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.statusManager != nil {
		return b.statusManager()
	}

	return BindStatusUnknown
}

func (b *BindBase) IsStatus(status BindStatus) bool {
	return b.Status() == status
}

func (b *BindBase) IsStatusUnknown() bool {
	return b.IsStatus(BindStatusUnknown)
}

func (b *BindBase) IsStatusUninitialized() bool {
	return b.IsStatus(BindStatusUninitialized)
}

func (b *BindBase) IsStatusInitializing() bool {
	return b.IsStatus(BindStatusInitializing)
}

func (b *BindBase) IsStatusOnline() bool {
	return b.IsStatus(BindStatusOnline)
}

func (b *BindBase) IsStatusOffline() bool {
	return b.IsStatus(BindStatusOffline)
}

func (b *BindBase) IsStatusRemoving() bool {
	return b.IsStatus(BindStatusRemoving)
}

func (b *BindBase) IsStatusRemoved() bool {
	return b.IsStatus(BindStatusRemoved)
}

func (b *BindBase) WrapTaskIsOnline(fn func(context.Context) error) *task.FunctionTask {
	return task.NewFunctionTask(func(ctx context.Context) (interface{}, error) {
		if !b.IsStatusOnline() {
			return nil, nil
		}

		return nil, fn(ctx)
	})
}

func (b *BindBase) SerialNumber() string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.serialNumber
}

func (b *BindBase) SetSerialNumber(serialNumber string) {
	b.mutex.Lock()
	b.serialNumber = serialNumber
	b.mutex.Unlock()
}

type BindMQTT struct {
	mutex  sync.RWMutex
	client mqtt.Component
}

func (b *BindMQTT) SetMQTTClient(client mqtt.Component) {
	b.mutex.Lock()
	b.client = client
	b.mutex.Unlock()
}

func (b *BindMQTT) MQTTClient() mqtt.Component {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.client
}

func (b *BindMQTT) MQTTPublish(ctx context.Context, topic mqtt.Topic, payload interface{}) error {
	return b.MQTTPublishRaw(ctx, topic, 1, true, payload)
}

func (b *BindMQTT) MQTTPublishWithoutCache(ctx context.Context, topic mqtt.Topic, payload interface{}) error {
	return b.MQTTPublishRawWithoutCache(ctx, topic, 1, true, payload)
}

func (b *BindMQTT) MQTTPublishRaw(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client := b.MQTTClient()

	if client == nil {
		return errors.New("MQTT client isn't init")
	}

	return client.Publish(ctx, topic, qos, retained, payload)
}

func (b *BindMQTT) MQTTPublishRawWithoutCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client := b.MQTTClient()

	if client == nil {
		return errors.New("MQTT client isn't init")
	}

	return client.PublishWithoutCache(ctx, topic, qos, retained, payload)
}

func (b *BindMQTT) MQTTPublishAsync(ctx context.Context, topic mqtt.Topic, payload interface{}) error {
	return b.MQTTPublishAsyncRaw(ctx, topic, 1, true, payload)
}

func (b *BindMQTT) MQTTPublishAsyncWithoutCache(ctx context.Context, topic mqtt.Topic, payload interface{}) error {
	return b.MQTTPublishAsyncRawWithoutCache(ctx, topic, 1, true, payload)
}

func (b *BindMQTT) MQTTPublishAsyncRaw(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client := b.MQTTClient()

	if client == nil {
		return errors.New("MQTT client isn't init")
	}

	client.PublishAsync(ctx, topic, qos, retained, payload)
	return nil
}

func (b *BindMQTT) MQTTPublishAsyncRawWithoutCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client := b.MQTTClient()

	if client == nil {
		return errors.New("MQTT client isn't init")
	}

	client.PublishAsyncWithoutCache(ctx, topic, qos, retained, payload)
	return nil
}

func CheckValueInMQTTTopic(topic mqtt.Topic, value string, offset int) bool {
	if value == "" {
		return false
	}

	routes := topic.Split()
	if len(routes) < offset {
		return false
	}

	return routes[len(routes)-offset] == value
}

func CheckSerialNumberInMQTTTopic(bind Bind, topic mqtt.Topic, offset int) bool {
	return CheckValueInMQTTTopic(topic, mqtt.NameReplace(bind.SerialNumber()), offset)
}

func WrapMQTTSubscribeDeviceIsOnline(status BindStatusManager, callback mqtt.MessageHandler) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
		if status() == BindStatusOnline {
			return callback(ctx, client, message)
		}

		return errors.New("bind isn't online")
	}
}
