package di

import (
	"context"
	"errors"
	"sync"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type MQTTHasSubscribers interface {
	MQTTSubscribers() []mqtt.Subscriber
}

type MQTTHasPublishes interface {
	MQTTPublishes() []mqtt.Topic
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

type MQTTContainer struct {
	bind boggart.BindItem

	clientMutex sync.RWMutex
	client      mqtt.Component

	cacheMutex       sync.Mutex
	cacheSubscribers []mqtt.Subscriber
	cachePublishes   []mqtt.Topic
}

func NewMQTTContainer(bind boggart.BindItem, client mqtt.Component) *MQTTContainer {
	return &MQTTContainer{
		bind:   bind,
		client: client,
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

	return client.Publish(ctx, topic, qos, retained, payload)
}

func (c *MQTTContainer) PublishRawWithoutCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client := c.Client()

	if client == nil {
		return errors.New("MQTT client isn't init")
	}

	return client.PublishWithoutCache(ctx, topic, qos, retained, payload)
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
	return nil
}

func (c *MQTTContainer) PublishAsyncRawWithoutCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client := c.Client()

	if client == nil {
		return errors.New("MQTT client isn't init")
	}

	client.PublishAsyncWithoutCache(ctx, topic, qos, retained, payload)
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

func (c *MQTTContainer) WrapSubscribeDeviceIsOnline(callback mqtt.MessageHandler) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
		if c.bind.Status() == boggart.BindStatusOnline {
			return callback(ctx, client, message)
		}

		return errors.New("bind isn't online")
	}
}

func (c *MQTTContainer) Subscribers() []mqtt.Subscriber {
	has, ok := c.bind.Bind().(MQTTHasSubscribers)
	if !ok {
		return nil
	}

	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	if c.cacheSubscribers == nil {
		c.cacheSubscribers = has.MQTTSubscribers()
	}

	return c.cacheSubscribers
}

func (c *MQTTContainer) Publishes() []mqtt.Topic {
	has, ok := c.bind.Bind().(MQTTHasPublishes)
	if !ok {
		return nil
	}

	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()

	if c.cachePublishes == nil {
		c.cachePublishes = has.MQTTPublishes()
	}

	return c.cachePublishes
}
