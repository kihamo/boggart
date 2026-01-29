package di

import (
	"context"
	"errors"
	"strconv"
	"sync"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/logging"
)

type BindHasMQTTSubscribers interface {
	MQTTSubscribers() []mqtt.Subscriber
}

type MQTTContainerSupport interface {
	SetMQTT(*MQTTContainer)
	MQTT() *MQTTContainer
}

func MQTTContainerBind(bind boggart.Bind) (*MQTTContainer, bool) {
	if support, ok := bind.(MQTTContainerSupport); ok {
		container := support.MQTT()
		return container, container != nil
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
	subscriber     *mqttWrapSubscriber
	success        *atomic.Uint64
	successMessage *atomic.Value
	failed         *atomic.Uint64
	failedMessage  *atomic.Value
}

func (s MQTTSubscriber) Subscriber() mqtt.Subscriber {
	return s.subscriber
}

func (s MQTTSubscriber) SuccessCount() uint64 {
	return s.success.Load()
}

func (s MQTTSubscriber) SuccessMessage() mqtt.Message {
	if v := s.successMessage.Load(); v != nil {
		return v.(mqtt.Message)
	}

	return nil
}

func (s MQTTSubscriber) FailedCount() uint64 {
	return s.failed.Load()
}

func (s MQTTSubscriber) FailedMessage() mqtt.Message {
	if v := s.failedMessage.Load(); v != nil {
		return v.(mqtt.Message)
	}

	return nil
}

type MQTTContainer struct {
	bindItem boggart.BindItem

	clientMutex sync.RWMutex
	client      mqtt.Component

	mutex       sync.RWMutex
	subscribers []MQTTSubscriber
	publishes   map[string]uint64
}

func NewMQTTContainer(bindItem boggart.BindItem, client mqtt.Component) *MQTTContainer {
	return &MQTTContainer{
		bindItem:  bindItem,
		client:    client,
		publishes: make(map[string]uint64),
	}
}

func (c *MQTTContainer) getClient() (mqtt.Component, error) {
	c.clientMutex.RLock()
	defer c.clientMutex.RUnlock()

	if c.client == nil {
		return nil, errors.New("MQTT client isn't init")
	}

	return c.client, nil
}

func (c *MQTTContainer) SetClient(client mqtt.Component) {
	c.clientMutex.Lock()
	c.client = client
	c.clientMutex.Unlock()
}

func (c *MQTTContainer) registerPublish(topic mqtt.Topic) {
	key := topic.String()

	c.mutex.Lock()
	c.publishes[key]++
	c.mutex.Unlock()
}

func (c *MQTTContainer) Publish(ctx context.Context, topic mqtt.Topic, payload interface{}) error {
	return c.PublishRaw(ctx, topic, 1, true, payload)
}

func (c *MQTTContainer) PublishWithoutCache(ctx context.Context, topic mqtt.Topic, payload interface{}) error {
	return c.PublishRawWithoutCache(ctx, topic, 1, true, payload)
}

func (c *MQTTContainer) PublishRaw(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	err = client.Publish(ctx, topic, qos, retained, payload)
	if err == nil {
		c.registerPublish(topic)
	}

	return err
}

func (c *MQTTContainer) PublishRawWithoutCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	err = client.PublishWithoutCache(ctx, topic, qos, retained, payload)
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
	client, err := c.getClient()
	if err != nil {
		return err
	}

	client.PublishAsync(ctx, topic, qos, retained, payload)
	c.registerPublish(topic)

	return nil
}

func (c *MQTTContainer) PublishAsyncRawWithoutCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	client.PublishAsyncWithoutCache(ctx, topic, qos, retained, payload)
	c.registerPublish(topic)

	return nil
}

func (c *MQTTContainer) Delete(ctx context.Context, topic mqtt.Topic) error {
	return c.PublishWithoutCache(ctx, topic, nil)
}

func (c *MQTTContainer) Request(ctx context.Context, requestTopic, responseTopic mqtt.Topic, requestPayload interface{}) (_ mqtt.Message, err error) {
	response := make(chan mqtt.Message, 1)
	single := make(chan struct{})

	subscribe := mqtt.NewSubscriber(responseTopic, 1, func(ctx context.Context, _ mqtt.Component, message mqtt.Message) error {
		// TODO: иногда рутина на этом месте блокируется навечно, возможно из-за того что больше одного сообщения приходит
		response <- message
		<-single

		return nil
	})

	err = c.Subscribe(subscribe)
	if err != nil {
		return nil, err
	}

	defer func() {
		close(single)
		c.Unsubscribe(subscribe)
	}()

	err = c.PublishRawWithoutCache(ctx, requestTopic, 1, false, requestPayload)
	if err != nil {
		return nil, err
	}

	select {
	case m := <-response:
		return m, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// смещение человеко понятное и стартует с единицы:
//
//	offset == 0 игнорируется
//	offset > 0 смещение от лева к праву 2 => (a [b] c d e f)
//	offset < 0 смещение от права к леву -2 => (a b c d [e] f)
func (c *MQTTContainer) CheckValueInTopic(topic mqtt.Topic, value string, offset int) bool {
	if value == "" || offset == 0 {
		return false
	}

	routes := topic.Split()

	if offset > 0 {
		if len(routes) < offset {
			return false
		}

		return routes[offset-1] == value
	}

	offset *= -1
	if len(routes) < offset {
		return false
	}

	return routes[len(routes)-offset] == value
}

func (c *MQTTContainer) CheckBindIDInTopic(topic mqtt.Topic, offset int) bool {
	if bindSupport, ok := MetaContainerBind(c.bindItem.Bind()); ok {
		if id := bindSupport.ID(); id != "" {
			return c.CheckValueInTopic(topic, mqtt.NameReplace(id), offset)
		}
	}

	return false
}

func (c *MQTTContainer) CheckSerialNumberInTopic(topic mqtt.Topic, offset int) bool {
	if bindSupport, ok := MetaContainerBind(c.bindItem.Bind()); ok {
		if sn := bindSupport.SerialNumber(); sn != "" {
			return c.CheckValueInTopic(topic, mqtt.NameReplace(sn), offset)
		}
	}

	return false
}

func (c *MQTTContainer) CheckMACInTopic(topic mqtt.Topic, offset int) bool {
	if bindSupport, ok := MetaContainerBind(c.bindItem.Bind()); ok {
		if mac := bindSupport.MACAsString(); mac != "" {
			return c.CheckValueInTopic(topic, mqtt.NameReplace(mac), offset)
		}
	}

	return false
}

func (c *MQTTContainer) WrapSubscribeDeviceIsOnline(callback mqtt.MessageHandler) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) error {
		if c.bindItem.Status().IsStatusOnline() {
			return callback(ctx, client, message)
		}

		return errors.New("bind isn't online")
	}
}

func (c *MQTTContainer) createSubscribe(subscriber mqtt.Subscriber) MQTTSubscriber {
	logger, _ := LoggerContainerBind(c.bindItem.Bind())

	item := MQTTSubscriber{
		success:        atomic.NewUint64(),
		successMessage: atomic.NewValue(),
		failed:         atomic.NewUint64(),
		failedMessage:  atomic.NewValue(),
	}
	item.subscriber = newMQTTWrapSubscriber(subscriber, logger, func(message mqtt.Message) {
		item.success.Inc()
		item.successMessage.Store(message)
	}, func(message mqtt.Message) {
		item.failed.Inc()
		item.failedMessage.Store(message)
	})

	return item
}

func (c *MQTTContainer) Subscribe(subscribers ...mqtt.Subscriber) error {
	if len(subscribers) == 0 {
		return nil
	}

	client, err := c.getClient()
	if err != nil {
		return err
	}

	for _, s := range subscribers {
		item := c.createSubscribe(s)

		err = client.SubscribeSubscriber(item.Subscriber())
		if err != nil {
			return err
		}

		c.mutex.Lock()
		c.subscribers = append(c.subscribers, item)
		c.mutex.Unlock()
	}

	return nil
}

func (c *MQTTContainer) Unsubscribe(subscribers ...mqtt.Subscriber) error {
	if len(subscribers) == 0 {
		return nil
	}

	client, err := c.getClient()
	if err != nil {
		return err
	}

	for _, s := range subscribers {
		err = client.UnsubscribeSubscriber(s)
		if err != nil {
			return err
		}

		c.mutex.Lock()

		for i := len(c.subscribers) - 1; i >= 0; i-- {
			// смотрим напрямую в подписчика, игнорируя обертку
			if c.subscribers[i].subscriber.original == s {
				c.subscribers = append(c.subscribers[:i], c.subscribers[i+1:]...)
			}
		}

		c.mutex.Unlock()
	}

	return nil
}

func (c *MQTTContainer) Subscribers() []MQTTSubscriber {
	c.mutex.RLock()
	if c.subscribers != nil {
		defer c.mutex.RUnlock()

		return append([]MQTTSubscriber(nil), c.subscribers...)
	}
	c.mutex.RUnlock()

	has, ok := c.bindItem.Bind().(BindHasMQTTSubscribers)
	if !ok {
		return nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, subscriber := range has.MQTTSubscribers() {
		c.subscribers = append(c.subscribers, c.createSubscribe(subscriber))
	}

	return c.subscribers
}

func (c *MQTTContainer) Publishes() map[mqtt.Topic]uint64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	list := make(map[mqtt.Topic]uint64, len(c.publishes))
	for topic, count := range c.publishes {
		list[mqtt.Topic(topic)] = count
	}

	return list
}

func (c *MQTTContainer) ClientOptions() (*m.ClientOptionsReader, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	opts := client.Client().OptionsReader()
	return &opts, nil
}

type mqttWrapSubscriber struct {
	original mqtt.Subscriber
	logger   logging.Logger
	success  func(mqtt.Message)
	failed   func(mqtt.Message)
}

func newMQTTWrapSubscriber(subscriber mqtt.Subscriber, logger logging.Logger, success, failed func(mqtt.Message)) *mqttWrapSubscriber {
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

		t.failed(message)
	} else {
		t.success(message)
	}

	return err
}
