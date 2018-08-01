package internal

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"time"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/mqtt/subscribes"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

type Component struct {
	application shadow.Application
	config      config.Component

	mutex  sync.RWMutex
	client m.Client
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
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}

func (c *Component) Run() (err error) {
	lg := logger.NewOrNop(c.Name(), c.application)

	m.ERROR = NewMQTTLogger(lg.Error, lg.Errorf)
	m.CRITICAL = NewMQTTLogger(lg.Panic, lg.Panicf)
	m.WARN = NewMQTTLogger(lg.Warn, lg.Warnf)
	m.DEBUG = NewMQTTLogger(lg.Debug, lg.Debugf)

	c.initClient()

	return nil
}

func (c *Component) initClient() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	servers := make([]*url.URL, 0)
	for _, u := range strings.Split(c.config.String(mqtt.ConfigServers), ";") {
		if p, err := url.Parse(u); err == nil {
			servers = append(servers, p)
		}
	}

	h := md5.New()
	io.WriteString(h, time.Now().Format(time.RFC3339Nano))

	opts := &m.ClientOptions{
		Servers:  servers,
		ClientID: fmt.Sprintf("%x\n", h.Sum(nil)),
		Username: c.config.String(mqtt.ConfigUsername),
		Password: c.config.String(mqtt.ConfigPassword),
	}

	client := m.NewClient(opts)
	client.Connect().Wait()

	c.client = client

	defaultSubscribers := []mqtt.Subscriber{
		subscribes.NewOwnTracksSubscribe(),
	}

	for _, s := range defaultSubscribers {
		client.SubscribeMultiple(s.Filters(), func(client m.Client, message m.Message) {
			s.Callback(c, message)
		})
	}
}

func (c *Component) Client() m.Client {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.client
}

func (c *Component) Publish(topic string, qos byte, retained bool, payload interface{}) m.Token {
	return c.Client().Publish(topic, qos, retained, payload)
}

func (c *Component) Subscribe(topic string, qos byte, callback m.MessageHandler) m.Token {
	return c.Client().Subscribe(topic, qos, callback)
}

func (c *Component) SubscribeMultiple(filters map[string]byte, callback m.MessageHandler) m.Token {
	return c.Client().SubscribeMultiple(filters, callback)
}

func (c *Component) Unsubscribe(topics ...string) m.Token {
	return c.Client().Unsubscribe(topics...)
}

func (c *Component) AddRoute(topic string, callback m.MessageHandler) {
	c.Client().AddRoute(topic, callback)
}

func (c *Component) OptionsReader() m.ClientOptionsReader {
	return c.Client().OptionsReader()
}
