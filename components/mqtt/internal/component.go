package internal

import (
	"bytes"
	"container/list"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
)

const timeFormat = "2006-01-02T15:04:05+0000"

type Component struct {
	lostConnections uint64

	application shadow.Application
	components  []shadow.Component
	config      config.Component
	logger      logging.Logger

	mutex             sync.RWMutex
	routes            []dashboard.Route
	client            m.Client
	subscriptions     *list.List
	handlersOnConnect []mqtt.OnConnectHandler
	payloadCache      mqtt.Cache
}

func (c *Component) Name() string {
	return mqtt.ComponentName
}

func (c *Component) Version() string {
	return mqtt.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: logging.ComponentName,
		},
		//{
		//	Name: tracing.ComponentName,
		//},
	}
}

func (c *Component) Init(a shadow.Application) (err error) {
	c.application = a

	c.subscriptions = list.New()
	c.handlersOnConnect = make([]mqtt.OnConnectHandler, 0)

	c.payloadCache, err = newCache(1)

	return err
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) (err error) {
	if c.components, err = a.GetComponents(); err != nil {
		return err
	}

	c.logger = logging.DefaultLazyLogger(c.Name())

	clientLogger := logging.NewLazyLogger(c.logger, c.Name()+".client")
	m.ERROR = NewMQTTLogger(clientLogger.Error, clientLogger.Errorf)
	m.CRITICAL = m.ERROR
	m.WARN = NewMQTTLogger(clientLogger.Warn, clientLogger.Warnf)
	m.DEBUG = NewMQTTLogger(clientLogger.Debug, clientLogger.Debugf)

	<-a.ReadyComponent(config.ComponentName)
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	if err = c.payloadCache.Resize(c.config.Int(mqtt.ConfigPayloadCacheSize)); err != nil {
		return err
	}

	if err := c.initClient(); err != nil {
		return err
	}

	ready <- struct{}{}

	return c.initSubscribers()
}

func (c *Component) initClient() (err error) {
	opts := m.NewClientOptions()
	opts.SetStore(NewStore(c.logger))

	opts.SetClientID(c.config.String(mqtt.ConfigClientID))
	opts.SetUsername(c.config.String(mqtt.ConfigUsername))
	opts.SetPassword(c.config.String(mqtt.ConfigPassword))
	opts.SetCleanSession(c.config.Bool(mqtt.ConfigClearSession))
	opts.SetResumeSubs(c.config.Bool(mqtt.ConfigResumeSubs))
	opts.SetOrderMatters(c.config.Bool(mqtt.ConfigOrderMatters))

	opts.SetConnectTimeout(c.config.Duration(mqtt.ConfigConnectionTimeout))
	opts.SetConnectRetry(c.config.Bool(mqtt.ConfigConnectionRetryEnabled))
	opts.SetConnectRetryInterval(c.config.Duration(mqtt.ConfigConnectionRetryInterval))

	opts.SetKeepAlive(c.config.Duration(mqtt.ConfigKeepAlive))
	opts.SetWriteTimeout(c.config.Duration(mqtt.ConfigWriteTimeout))
	opts.SetPingTimeout(c.config.Duration(mqtt.ConfigPingTimeout))

	opts.SetAutoReconnect(c.config.Bool(mqtt.ConfigReconnectEnabled))
	opts.SetMaxReconnectInterval(c.config.Duration(mqtt.ConfigReconnectMaxInterval))

	if c.config.Bool(mqtt.ConfigLWTEnabled) {
		opts.SetWill(
			c.config.String(mqtt.ConfigLWTTopic),
			c.config.String(mqtt.ConfigLWTPayload),
			byte(c.config.Int(mqtt.ConfigLWTQOS)),
			c.config.Bool(mqtt.ConfigLWTRetained))
	} else {
		opts.UnsetWill()
	}

	// specific handlers
	opts.SetOnConnectHandler(func(client m.Client) {
		cfg := client.OptionsReader()

		var mqttVersion string

		switch cfg.ProtocolVersion() {
		case 3:
			mqttVersion = "3.1"
		case 4:
			mqttVersion = "3.1.1"
		}

		c.logger.Debug("Connect to MQTT broker", "clientId", cfg.ClientID(), "protocol.version", mqttVersion)
		metricConnect.Inc()

		restore := atomic.LoadUint64(&c.lostConnections) == 0

		if !restore {
			for _, sub := range c.Subscriptions() {
				topic := sub.Topic()
				qos := sub.QOS()

				if err := c.clientSubscribe(topic, qos, sub); err != nil {
					c.logger.Error("Resubscribe failed", "topic", topic, "qos", qos, "error", err.Error())
				} else {
					c.logger.Debug("Resubscribe success", "topic", topic, "qos", qos)
				}
			}
		}

		c.mutex.RLock()
		for _, handler := range c.handlersOnConnect {
			go handler(c, !restore)
		}
		c.mutex.RUnlock()
	})

	opts.SetConnectionLostHandler(func(client m.Client, reason error) {
		atomic.AddUint64(&c.lostConnections, 1)
		c.logger.Error("Connection lost", "error", reason.Error(), "count", atomic.LoadUint64(&c.lostConnections))
		metricConnectionLost.Inc()
	})

	opts.SetReconnectingHandler(func(_ m.Client, _ *m.ClientOptions) {
		c.logger.Debug("Attempt reconnect")
	})

	opts.SetDefaultPublishHandler(func(_ m.Client, message m.Message) {
		c.logger.Warn("Received that does not match any known subscriptions",
			"topic", message.Topic(),
			"qos", message.Qos(),
			"retained", message.Retained(),
		)
	})

	for _, u := range strings.Split(c.config.String(mqtt.ConfigServers), ";") {
		opts.AddBroker(u)
	}

	client := m.NewClient(opts)
	err = c.tokenWait(client.Connect())

	c.mutex.Lock()
	c.client = client
	c.mutex.Unlock()

	return err
}

func (c *Component) initSubscribers() error {
	for _, component := range c.components {
		if subscribers, ok := component.(mqtt.HasSubscribers); ok {
			for _, subscriber := range subscribers.MQTTSubscribers() {
				if err := c.SubscribeSubscriber(subscriber); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *Component) tokenWait(token m.Token) error {
	timeout := c.config.Duration(mqtt.ConfigTokenWaitTimeout)

	if !token.WaitTimeout(timeout) {
		return errors.New("token wait didn't return in an expected time " + timeout.String())
	}

	return token.Error()
}

func (c *Component) Shutdown() (err error) {
	c.mutex.RLock()
	client := c.client
	c.mutex.RUnlock()

	if client != nil {
		if c.config.Bool(mqtt.ConfigLWTEnabled) {
			err = c.tokenWait(client.Publish(
				c.config.String(mqtt.ConfigLWTTopic),
				byte(c.config.Int(mqtt.ConfigLWTQOS)),
				c.config.Bool(mqtt.ConfigLWTRetained),
				[]byte(c.config.String(mqtt.ConfigLWTPayload))))
		}

		client.Disconnect(250)
	}

	return err
}

func (c *Component) Client() m.Client {
	c.mutex.RLock()
	client := c.client
	c.mutex.RUnlock()

	if client != nil {
		return client
	}

	<-c.application.ReadyComponent(c.Name())

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.client
}

func (c *Component) clientSubscribe(topic mqtt.Topic, qos byte, subscription *mqtt.Subscription) error {
	client := c.Client()
	if client == nil {
		return errors.New("can't initialize client of MQTT")
	}

	// check topic
	if !topic.ValidAsSubscribeTopic() {
		return errors.New("topic " + topic.String() + " isn't valid for subscribe")
	}

	// wrap tracing
	callback := func(client m.Client, message m.Message) {
		ctx := context.Background()
		//span, ctx := tracing.StartSpanFromContext(context.Background(), c.Name(), "subscribe_callback")

		//span = span.SetTag("topic", message.Topic())
		//defer span.Finish()
		//
		//span.LogFields(
		//	log.Int("qos", int(message.Qos())),
		//	log.String("payload", string(message.Payload())),
		//	log.Bool("retained", message.Retained()),
		//	log.String("topic.subscribe", topic.String()),
		//)

		msg := newMessage(message)

		// в отдельной рутине, так как если зависнет хендлер клиент MQTT не сделает ack на сообщение
		go func() {
			logPayload := msg.String()
			if len(logPayload) > 100 {
				logPayload = logPayload[:100]
			}

			c.logger.Debug(
				"Call MQTT subscriber",
				"topic.subscribe", topic,
				"topic.call", message.Topic(),
				"qos", strconv.Itoa(int(qos)),
				"retained", strconv.FormatBool(message.Retained()),
				"payload", logPayload,
			)

			if err := subscription.Callback(ctx, c, msg); err != nil {
				metricSubscriberCalls.With("status", "failure", "topic", topic.String()).Inc()

				//tracing.SpanError(span, err)

				c.logger.Error(
					"Call MQTT subscriber failed",
					"error", err.Error(),
					"topic.subscribe", topic,
					"topic.call", message.Topic(),
					"qos", strconv.Itoa(int(qos)),
					"retained", strconv.FormatBool(message.Retained()),
					"payload", logPayload,
				)
			} else {
				metricSubscriberCalls.With("status", "success", "topic", topic.String()).Inc()

				c.logger.Debug(
					"Call MQTT subscriber success",
					"topic.subscribe", topic,
					"topic.call", message.Topic(),
					"qos", strconv.Itoa(int(qos)),
					"retained", strconv.FormatBool(message.Retained()),
					"payload", logPayload,
				)
			}
		}()
	}

	token := client.Subscribe(topic.String(), qos, callback)

	err := c.tokenWait(token)
	if err == nil {
		metricSubscribe.With("status", "success", "topic", topic.String()).Inc()
	} else {
		metricSubscribe.With("status", "failure", "topic", topic.String()).Inc()
	}

	return err
}

func (c *Component) doPublish(_ context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}, cache bool) (err error) {
	// check topic
	if !topic.ValidAsPublishTopic() {
		return errors.New("topic " + topic.String() + " isn't valid for publish")
	}

	payloadConverted := c.convertPayload(payload)

	if cache && !c.config.Bool(mqtt.ConfigPayloadCacheEnabled) {
		cache = false
	}

	var logPayload string
	if len(payloadConverted) > 100 {
		logPayload = string(payloadConverted[:100])
	} else {
		logPayload = string(payloadConverted)
	}

	logQOS := strconv.Itoa(int(qos))
	logRetained := strconv.FormatBool(retained)

	client := c.Client()
	if client != nil {
		if cache {
			if val, ok := c.payloadCache.Get(topic); ok && bytes.Equal(val.Payload(), payloadConverted) {
				metricPayloadCacheHit.Inc()

				c.logger.Debug(
					"Publish MQTT topic skip because there is cache hit",
					"topic", topic,
					"qos", logQOS,
					"retained", logRetained,
					"payload", logPayload,
				)

				return nil
			}

			metricPayloadCacheMiss.Inc()
		}

		err = c.tokenWait(client.Publish(topic.String(), qos, retained, payloadConverted))
	} else {
		err = errors.New("can't initialize client of MQTT")
	}

	if err != nil {
		metricPublish.With("status", "failure").Inc()

		c.logger.Error(
			"Publish MQTT topic failed",
			"error", err.Error(),
			"topic", topic,
			"qos", logQOS,
			"retained", logRetained,
			"payload", logPayload,
		)
	} else {
		metricPublish.With("status", "success").Inc()

		// if cache {
		c.payloadCache.Add(topic, payloadConverted)
		// }
	}

	return err
}

func (c *Component) Publish(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	return c.doPublish(ctx, topic, qos, retained, payload, retained)
}

func (c *Component) PublishWithCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	return c.doPublish(ctx, topic, qos, retained, payload, true)
}

func (c *Component) PublishWithoutCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) error {
	return c.doPublish(ctx, topic, qos, retained, payload, false)
}

func (c *Component) PublishAsync(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) {
	go func() {
		_ = c.doPublish(ctx, topic, qos, retained, payload, retained)
	}()
}

func (c *Component) PublishAsyncWithCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) {
	go func() {
		_ = c.doPublish(ctx, topic, qos, retained, payload, true)
	}()
}

func (c *Component) PublishAsyncWithoutCache(ctx context.Context, topic mqtt.Topic, qos byte, retained bool, payload interface{}) {
	go func() {
		_ = c.doPublish(ctx, topic, qos, retained, payload, false)
	}()
}

func (c *Component) Unsubscribe(topic mqtt.Topic) error {
	client := c.Client()
	if client == nil {
		return errors.New("can't initialize client of MQTT")
	}

	if err := c.tokenWait(client.Unsubscribe(topic.String())); err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for element := c.subscriptions.Front(); element != nil; element = element.Next() {
		subscription := element.Value.(*mqtt.Subscription)

		for _, subscriber := range subscription.Subscribers() {
			if subscriber.Topic().String() == topic.String() {
				subscription.RemoveSubscriber(subscriber)
			}
		}

		if subscription.Len() == 0 {
			topic := subscription.Topic().String()

			if err := c.tokenWait(client.Unsubscribe(topic)); err == nil {
				c.subscriptions.Remove(element)
				c.logger.Debug("Unsubscribe", "topic", topic)
			} else {
				c.logger.Error("Unsubscribe failed", "error", err.Error())
				return err
			}
		}
	}

	return nil
}

func (c *Component) UnsubscribeSubscriber(subscriber mqtt.Subscriber) error {
	client := c.Client()
	if client == nil {
		return errors.New("can't initialize client of MQTT")
	}

	var element *list.Element

	c.mutex.Lock()
	defer c.mutex.Unlock()

	for element = c.subscriptions.Front(); element != nil; element = element.Next() {
		subscription := element.Value.(*mqtt.Subscription)

		if result := subscription.RemoveSubscriber(subscriber); result {
			if subscription.Len() == 0 {
				topic := subscriber.Topic().String()

				if err := c.tokenWait(client.Unsubscribe(topic)); err == nil {
					c.subscriptions.Remove(element)
					c.logger.Debug("Unsubscribe", "topic", topic)
				} else {
					c.logger.Error("Unsubscribe failed", "error", err.Error())

					// fallback
					subscription.AddSubscriber(subscriber)

					return err
				}
			}

			break
		}
	}

	return nil
}

func (c *Component) UnsubscribeSubscribers(subscribers []mqtt.Subscriber) error {
	for _, subscriber := range subscribers {
		if err := c.UnsubscribeSubscriber(subscriber); err != nil {
			return err
		}
	}

	return nil
}

func (c *Component) Subscribe(topic mqtt.Topic, qos byte, callback mqtt.MessageHandler) (mqtt.Subscriber, error) {
	subscriber := mqtt.NewSubscriber(topic, qos, callback)
	if err := c.SubscribeSubscriber(subscriber); err != nil {
		return nil, err
	}

	return subscriber, nil
}

func (c *Component) SubscribeSubscriber(subscriber mqtt.Subscriber) error {
	var (
		element      *list.Element
		subscription *mqtt.Subscription
	)

	topic := subscriber.Topic()
	qos := subscriber.QOS()

	//  ищем подходящие подписки
	c.mutex.RLock()
	for element = c.subscriptions.Front(); element != nil; element = element.Next() {
		s := element.Value.(*mqtt.Subscription)
		if s.Match(topic) {
			s.AddSubscriber(subscriber)
			subscription = s

			break
		}
	}
	c.mutex.RUnlock()

	// если подписка не найдена, то создаем новую
	if subscription == nil {
		subscription = mqtt.NewSubscription(subscriber)
	}

	if err := c.clientSubscribe(topic, qos, subscription); err != nil {
		return err
	}

	c.mutex.Lock()
	if element != nil {
		element.Value = subscription
	} else {
		c.subscriptions.PushBack(subscription)
	}
	c.mutex.Unlock()

	c.logger.Debug("Subscribe success", "topic", topic, "qos", qos)

	return nil
}

func (c *Component) Subscriptions() []*mqtt.Subscription {
	c.mutex.RLock()
	subscriptions := make([]*mqtt.Subscription, 0, c.subscriptions.Len())

	for e := c.subscriptions.Front(); e != nil; e = e.Next() {
		subscriptions = append(subscriptions, e.Value.(*mqtt.Subscription))
	}

	c.mutex.RUnlock()

	return subscriptions
}

func (c *Component) OnConnectHandlerAdd(handler mqtt.OnConnectHandler) {
	c.mutex.Lock()
	c.handlersOnConnect = append(c.handlersOnConnect, handler)
	c.mutex.Unlock()
}

func (c *Component) CacheItems() []mqtt.CacheItem {
	return c.payloadCache.Payloads()
}

func (c *Component) convertPayload(payload interface{}) []byte {
	switch value := payload.(type) {
	case nil:
		return nil
	case []byte:
		return value
	case string:
		return []byte(value)
	case float64:
		return []byte(strconv.FormatFloat(value, 'f', -1, 64))
	case float32:
		return []byte(strconv.FormatFloat(float64(value), 'f', -1, 64))
	case int64:
		return []byte(strconv.FormatInt(value, 10))
	case int32:
		return []byte(strconv.FormatInt(int64(value), 10))
	case int16:
		return []byte(strconv.FormatInt(int64(value), 10))
	case int8:
		return []byte(strconv.FormatInt(int64(value), 10))
	case int:
		return []byte(strconv.FormatInt(int64(value), 10))
	case uint64:
		return []byte(strconv.FormatUint(value, 10))
	case uint32:
		return []byte(strconv.FormatUint(uint64(value), 10))
	case uint16:
		return []byte(strconv.FormatUint(uint64(value), 10))
	case uint8:
		return []byte(strconv.FormatUint(uint64(value), 10))
	case uint:
		return []byte(strconv.FormatUint(uint64(value), 10))
	case bool:
		if value {
			return PayloadTrue
		}

		return PayloadFalse
	case time.Time:
		return []byte(value.UTC().Format(timeFormat))
	case *time.Time:
		return []byte(value.UTC().Format(timeFormat))
	case time.Duration:
		return []byte(strconv.FormatFloat(value.Seconds(), 'f', -1, 64))
	case *time.Duration:
		return []byte(strconv.FormatFloat(value.Seconds(), 'f', -1, 64))
	case bytes.Buffer:
		return value.Bytes()
	case io.Reader:
		if b, err := ioutil.ReadAll(value); err == nil {
			return b
		}

		return []byte(fmt.Sprintf("%v", payload))
	case fmt.Stringer:
		return []byte(value.String())
	default:
		if ref := reflect.ValueOf(value); ref.Kind() == reflect.Ptr {
			if !ref.Elem().IsValid() {
				return c.convertPayload(nil)
			}

			return c.convertPayload(ref.Elem().Interface())
		}

		return []byte(fmt.Sprintf("%v", payload))
	}
}
