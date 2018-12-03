package internal

import (
	"container/list"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logging"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go/log"
)

type Component struct {
	application shadow.Application
	components  []shadow.Component
	config      config.Component
	logger      logging.Logger

	mutex         sync.RWMutex
	routes        []dashboard.Route
	client        m.Client
	subscriptions *list.List
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
		{
			Name: tracing.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) (err error) {
	if c.components, err = a.GetComponents(); err != nil {
		return err
	}

	c.subscriptions = list.New()
	c.logger = logging.DefaultLogger().Named(c.Name())

	m.ERROR = NewMQTTLogger(c.logger.Error, c.logger.Errorf)
	m.CRITICAL = NewMQTTLogger(c.logger.Panic, c.logger.Panicf)
	m.WARN = NewMQTTLogger(c.logger.Warn, c.logger.Warnf)
	m.DEBUG = NewMQTTLogger(c.logger.Debug, c.logger.Debugf)

	<-a.ReadyComponent(config.ComponentName)
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	if err := c.initClient(); err != nil {
		return err
	}

	ready <- struct{}{}

	return c.initSubscribers()
}

func (c *Component) initClient() error {
	attempts := c.config.Int64(mqtt.ConfigConnectionAttempts)

	// auto reconnect
	//duration := c.config.Duration(mqtt.ConfigConnectionTimeout) + time.Second*30
	duration := time.Second * 5
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		opts := m.NewClientOptions()
		opts.Username = c.config.String(mqtt.ConfigUsername)
		opts.Password = c.config.String(mqtt.ConfigPassword)
		opts.ConnectTimeout = c.config.Duration(mqtt.ConfigConnectionTimeout)

		opts.Servers = make([]*url.URL, 0)
		for _, u := range strings.Split(c.config.String(mqtt.ConfigServers), ";") {
			if p, err := url.Parse(u); err == nil {
				opts.Servers = append(opts.Servers, p)
			}
		}

		h := md5.New()
		io.WriteString(h, time.Now().Format(time.RFC3339Nano))
		opts.ClientID = fmt.Sprintf("%x\n", h.Sum(nil))

		client := m.NewClient(opts)
		token := client.Connect()

		token.Wait()
		attempts--

		if err := token.Error(); err != nil {
			c.logger.Error("Init MQTT client failed with error " + err.Error())

			if attempts == 0 {
				c.logger.Error("Init MQTT client failed because exhausted number of connection attempts")
				return errors.New("exhausted number of connection attempts")
			}
		} else {
			c.mutex.Lock()
			c.client = client
			c.mutex.Unlock()

			break
		}
	}

	return nil
}

func (c *Component) initSubscribers() error {
	for _, component := range c.components {
		if componentSubscribers, ok := component.(mqtt.HasSubscribers); ok {
			for _, sub := range componentSubscribers.MQTTSubscribers() {
				if err := c.Subscribe(sub.Topic(), sub.QOS(), sub.Callback()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *Component) Shutdown() error {
	c.mutex.RLock()
	client := c.client
	c.mutex.RUnlock()

	if client != nil {
		client.Disconnect(250)
	}

	return nil
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

func (c *Component) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	r := "0"
	if retained {
		r = "1"
	}

	metricPublish.With(
		"topic", topic,
		"qos", strconv.Itoa(int(qos)),
		"retained", r,
	).Inc()

	token := c.Client().Publish(topic, qos, retained, payload)
	token.Wait()

	return token.Error()
}

func (c *Component) AddRoute(topic string, callback mqtt.MessageHandler) {
	subscription := mqtt.NewSubscription(topic, 0, callback)

	var e *list.Element
	c.mutex.RLock()
	for e = c.subscriptions.Front(); e != nil; e = e.Next() {
		if e.Value.(mqtt.Subscription).Match(topic) {
			subscription = subscription.Merge(e.Value.(mqtt.Subscription))
			break
		}
	}
	c.mutex.RUnlock()

	c.Client().AddRoute(topic, c.wrapCallback(subscription.Callback))

	c.mutex.Lock()
	if e != nil {
		e.Value = subscription
	} else {
		c.subscriptions.PushBack(subscription)
	}
	c.mutex.Unlock()
}

func (c *Component) Unsubscribe(topic string) error {
	if token := c.Client().Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	c.mutex.Lock()
	for e := c.subscriptions.Front(); e != nil; e = e.Next() {
		if e.Value.(mqtt.Subscription).Match(topic) {
			c.subscriptions.Remove(e)
			break
		}
	}
	c.mutex.Unlock()

	c.logger.Debug("Unsubscribe", "topic", topic)

	return nil
}

func (c *Component) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	subscription := mqtt.NewSubscription(topic, qos, callback)

	var e *list.Element
	c.mutex.RLock()
	for e = c.subscriptions.Front(); e != nil; e = e.Next() {
		if e.Value.(mqtt.Subscription).Match(topic) {
			subscription = subscription.Merge(e.Value.(mqtt.Subscription))
			break
		}
	}
	c.mutex.RUnlock()

	token := c.Client().Subscribe(topic, qos, c.wrapCallback(subscription.Callback))
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	c.mutex.Lock()
	if e != nil {
		e.Value = subscription
	} else {
		c.subscriptions.PushBack(subscription)
	}
	c.mutex.Unlock()

	c.logger.Debug("Subscribe", "topic", topic, "qos", qos)

	return nil
}

func (c *Component) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) error {
	for topic, qos := range filters {
		if err := c.Subscribe(topic, qos, callback); err != nil {
			return err
		}
	}

	return nil
}

func (c *Component) SubscribeSubscribers(subscribers []mqtt.Subscriber) error {
	for _, sub := range subscribers {
		if err := c.Subscribe(sub.Topic(), sub.QOS(), sub.Callback()); err != nil {
			return err
		}
	}

	return nil
}

func (c *Component) Subscriptions() []mqtt.Subscription {
	c.mutex.RLock()
	subscriptions := make([]mqtt.Subscription, 0, c.subscriptions.Len())
	for e := c.subscriptions.Front(); e != nil; e = e.Next() {
		subscriptions = append(subscriptions, e.Value.(mqtt.Subscription))
	}
	c.mutex.RUnlock()

	return subscriptions
}

func (c *Component) wrapCallback(callback mqtt.MessageHandler) func(client m.Client, message m.Message) {
	return func(client m.Client, message m.Message) {
		span, ctx := tracing.StartSpanFromContext(context.Background(), c.Name(), message.Topic())
		defer span.Finish()

		span.LogFields(log.String("payload", string(message.Payload())))

		r := "0"
		if message.Retained() {
			r = "1"
		}

		metricSubscribe.With(
			"topic", message.Topic(),
			"qos", strconv.Itoa(int(message.Qos())),
			"retained", r,
		).Inc()

		callback(ctx, c, message)
	}
}
