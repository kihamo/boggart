package di

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/logging"
)

type MQTTHasSubscribers interface {
	MQTTSubscribers() []mqtt.Subscriber
}

type MQTTContainerSupport interface {
	SetMQTT(*MQTTContainer)
	MQTT() *MQTTContainer
}

func MQTTContainerBind(bind boggart.Bind) (*MQTTContainer, bool) {
	if support, ok := bind.(MQTTContainerSupport); ok {
		return support.MQTT(), true
	}

	return nil, false
}

type MQTTBind struct {
	mutex     sync.RWMutex
	container *MQTTContainer
}

func (b *MQTTBind) SetMQTT(container *MQTTContainer) {
	b.mutex.Lock()
	b.container = container
	b.mutex.Unlock()
}

func (b *MQTTBind) MQTT() *MQTTContainer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.container
}

type MQTTSubscriber struct {
	subscriber  mqtt.Subscriber
	success     *atomic.Uint64
	successTime *atomic.TimeNull
	failed      *atomic.Uint64
	failedTime  *atomic.TimeNull
}

func (s MQTTSubscriber) Subscriber() mqtt.Subscriber {
	return s.subscriber
}

func (s MQTTSubscriber) SuccessCount() uint64 {
	return s.success.Load()
}

func (s MQTTSubscriber) SuccessTime() *time.Time {
	return s.successTime.Load()
}

func (s MQTTSubscriber) FailedCount() uint64 {
	return s.failed.Load()
}

func (s MQTTSubscriber) FailedTime() *time.Time {
	return s.failedTime.Load()
}

type MQTTContainer struct {
	bind boggart.BindItem

	clientMutex sync.RWMutex
	client      mqtt.Component

	cacheMutex       sync.RWMutex
	cacheSubscribers []MQTTSubscriber
	cachePublishes   map[string]uint64
}

func NewMQTTContainer(bind boggart.BindItem, client mqtt.Component) *MQTTContainer {
	return &MQTTContainer{
		bind:           bind,
		client:         client,
		cachePublishes: make(map[string]uint64, 0),
	}
}

func (c *MQTTContainer) Client() mqtt.Component {
	c.clientMutex.RLock()
	defer c.clientMutex.RUnlock()

	return c.client
}

func (c *MQTTContainer) SetClient(client mqtt.Component) {
	c.clientMutex.Lock()
	c.client = client
	c.clientMutex.Unlock()
}

func (c *MQTTContainer) registerPublish(topic mqtt.Topic) {
	key := topic.String()

	c.cacheMutex.Lock()
	if _, ok := c.cachePublishes[key]; ok {
		c.cachePublishes[key]++
	} else {
		c.cachePublishes[key] = 1
	}
	c.cacheMutex.Unlock()
}

func (c *MQTTContainer) Publish(ctx context.Context, topic mqtt.Topic, payload interface{}) error {
	return c.PublishRaw(ctx, topic, 1, true, payload)
}

func (c *MQTTContainer) PublishWithoutCache(ctx context.Context, topic mqtt.Topic, payload interface{}) error {
	return c.PublishRawWithoutCache(ctx, topic, 1, true, payload)
}

func (c *MQTTContainer) PublishRaw(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client := c.Client()

	if client == nil {
		return errors.New("MQTT client isn't init")
	}

	err := client.Publish(ctx, topic, qos, retained, payload)
	if err == nil {
		c.registerPublish(topic)
	}

	return err
}

func (c *MQTTContainer) PublishRawWithoutCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client := c.Client()

	if client == nil {
		return errors.New("MQTT client isn't init")
	}

	err := client.PublishWithoutCache(ctx, topic, qos, retained, payload)
	if err == nil {
		c.registerPublish(topic)
	}

	return err
}

func (c *MQTTContainer) PublishAsync(ctx context.Context, topic mqtt.Topic, payload interface{}) error {
	return c.PublishAsyncRaw(ctx, topic, 1, true, payload)
}

func (c *MQTTContainer) PublishAsyncWithoutCache(ctx context.Context, topic mqtt.Topic, payload interface{}) error {
	return c.PublishAsyncRawWithoutCache(ctx, topic, 1, true, payload)
}

func (c *MQTTContainer) PublishAsyncRaw(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client := c.Client()

	if client == nil {
		return errors.New("MQTT client isn't init")
	}

	client.PublishAsync(ctx, topic, qos, retained, payload)
	c.registerPublish(topic)

	return nil
}

func (c *MQTTContainer) PublishAsyncRawWithoutCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client := c.Client()

	if client == nil {
		return errors.New("MQTT client isn't init")
	}

	client.PublishAsyncWithoutCache(ctx, topic, qos, retained, payload)
	c.registerPublish(topic)

	return nil
}

func (c *MQTTContainer) CheckValueInTopic(topic mqtt.Topic, value string, offset int) bool {
	if value == "" {
		return false
	}

	routes := topic.Split()
	if len(routes) < offset {
		return false
	}

	return routes[len(routes)-offset] == value
}

func (c *MQTTContainer) CheckSerialNumberInTopic(topic mqtt.Topic, offset int) bool {
	if bindSupport, ok := MetaContainerBind(c.bind.Bind()); ok {
		if sn := bindSupport.SerialNumber(); sn != "" {
			return c.CheckValueInTopic(topic, mqtt.NameReplace(sn), offset)
		}
	}

	return false
}

func (c *MQTTContainer) CheckMACInTopic(topic mqtt.Topic, offset int) bool {
	if bindSupport, ok := MetaContainerBind(c.bind.Bind()); ok {
		if mac := bindSupport.MACAsString(); mac != "" {
			return c.CheckValueInTopic(topic, mqtt.NameReplace(mac), offset)
		}
	}

	return false
}

func (c *MQTTContainer) WrapSubscribeDeviceIsOnline(callback mqtt.MessageHandler) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
		if c.bind.Status() == boggart.BindStatusOnline {
			return callback(ctx, client, message)
		}

		return errors.New("bind isn't online")
	}
}

func (c *MQTTContainer) Subscribers() []MQTTSubscriber {
	has, ok := c.bind.Bind().(MQTTHasSubscribers)
	if !ok {
		return nil
	}

	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	if c.cacheSubscribers == nil {
		var logger logging.Logger

		if bindSupport, ok := c.bind.Bind().(LoggerContainerSupport); ok {
			logger = bindSupport.Logger()
		}

		for _, subscriber := range has.MQTTSubscribers() {
			item := MQTTSubscriber{
				success:     atomic.NewUint64(),
				successTime: atomic.NewTimeNull(),
				failed:      atomic.NewUint64(),
				failedTime:  atomic.NewTimeNull(),
			}
			item.subscriber = newMQTTWrapSubscriber(subscriber, logger, func() {
				item.success.Inc()
				item.successTime.Set(time.Now())
			}, func() {
				item.failed.Inc()
				item.failedTime.Set(time.Now())
			})

			c.cacheSubscribers = append(c.cacheSubscribers, item)
		}
	}

	return c.cacheSubscribers
}

func (c *MQTTContainer) Publishes() map[mqtt.Topic]uint64 {
	c.cacheMutex.RLock()
	defer c.cacheMutex.RUnlock()

	list := make(map[mqtt.Topic]uint64, len(c.cachePublishes))
	for topic, count := range c.cachePublishes {
		list[mqtt.Topic(topic)] = count
	}

	return list
}

type mqttWrapSubscriber struct {
	original mqtt.Subscriber
	logger   logging.Logger
	success  func()
	failed   func()
}

func newMQTTWrapSubscriber(subscriber mqtt.Subscriber, logger logging.Logger, success, failed func()) *mqttWrapSubscriber {
	return &mqttWrapSubscriber{
		original: subscriber,
		logger:   logger,
		success:  success,
		failed:   failed,
	}
}

func (t *mqttWrapSubscriber) Topic() mqtt.Topic {
	return t.original.Topic()
}

func (t *mqttWrapSubscriber) QOS() byte {
	return t.original.QOS()
}

func (t *mqttWrapSubscriber) Call(ctx context.Context, client mqtt.Component, message mqtt.Message) (err error) {
	err = t.original.Call(ctx, client, message)
	if err != nil {
		if t.logger != nil {
			logPayload := message.String()
			if len(logPayload) > 100 {
				logPayload = logPayload[:100]
			}

			t.logger.Error(
				"Call MQTT subscriber failed",
				"error", err.Error(),
				"topic.subscribe", t.Topic(),
				"topic.call", message.Topic(),
				"qos", strconv.Itoa(int(t.QOS())),
				"retained", strconv.FormatBool(message.Retained()),
				"payload", logPayload,
			)
		}

		t.failed()
	} else {
		t.success()
	}

	return err
}
