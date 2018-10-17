package internal

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
	tracer      opentracing.Tracer

	mutex       sync.RWMutex
	client      m.Client
	subscribers []mqtt.Subscriber
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
			Name: logger.ComponentName,
		},
		{
			Name: tracing.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.subscribers = make([]mqtt.Subscriber, 0)

	return nil
}

func (c *Component) Run() (err error) {
	c.logger = logger.NewOrNop(c.Name(), c.application)
	c.tracer = tracing.NewOrNop(c.application)

	m.ERROR = NewMQTTLogger(c.logger.Error, c.logger.Errorf)
	m.CRITICAL = NewMQTTLogger(c.logger.Panic, c.logger.Panicf)
	m.WARN = NewMQTTLogger(c.logger.Warn, c.logger.Warnf)
	m.DEBUG = NewMQTTLogger(c.logger.Debug, c.logger.Debugf)

	// auto reconnect
	duration := c.config.Duration(mqtt.ConfigConnectionTimeout) + time.Second*30
	ticker := time.NewTicker(duration)

	for ; true; <-ticker.C {
		err := c.initClient()
		if err != nil {
			c.logger.Errorf("Init MQTT client failed with error %s", err.Error())
		} else {
			break
		}
	}

	ticker.Stop()

	return nil
}

func (c *Component) Shutdown() error {
	client := c.Client()

	if client != nil {
		client.Disconnect(250)
	}

	return nil
}

func (c *Component) initClient() (err error) {
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
	err = token.Error()
	if err == nil {
		c.mutex.Lock()
		c.client = client

		for _, subscriber := range c.subscribers {
			c.subscribeByClient(client, subscriber)
		}
		c.mutex.Unlock()
	}

	return err
}

func (c *Component) Client() m.Client {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.client
}

func (c *Component) Publish(topic string, qos byte, retained bool, payload interface{}) m.Token {
	client := c.Client()
	if client == nil {
		c.logger.Warn("Client isn't init. Publish skipping", map[string]interface{}{
			"topic":    topic,
			"qos":      qos,
			"retained": retained,
			"payload":  payload,
		})

		return nil
	}

	return client.Publish(topic, qos, retained, payload)
}

func (c *Component) Subscribe(subscriber mqtt.Subscriber) m.Token {
	c.mutex.Lock()
	c.subscribers = append(c.subscribers, subscriber)
	c.mutex.Unlock()

	client := c.Client()
	if client == nil {
		return nil
	}

	return c.subscribeByClient(client, subscriber)
}

func (c *Component) subscribeByClient(client m.Client, subscriber mqtt.Subscriber) m.Token {
	c.logger.Debug("Add subscriber", map[string]interface{}{
		"filters": subscriber.Filters(),
	})

	return client.SubscribeMultiple(subscriber.Filters(), func(client m.Client, message m.Message) {
		span := c.tracer.StartSpan(message.Topic())
		ext.Component.Set(span, c.Name())
		defer span.Finish()

		subscriber.Callback(opentracing.ContextWithSpan(context.Background(), span), c, message)
	})
}
