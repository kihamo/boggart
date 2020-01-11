package components

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
)

type Base struct {
	ID                  string        `json:"-"`
	Type                ComponentType `json:"-"`
	Icon                string        `json:"icon"`
	Name                string        `json:"name"`
	StateTopic          mqtt.Topic    `json:"state_topic"`
	CommandTopic        mqtt.Topic    `json:"command_topic"`
	PayloadAvailable    string        `json:"payload_available"`
	PayloadNotAvailable string        `json:"payload_not_available"`
	UniqueID            string        `json:"unique_id"`
	Device              struct {
		Identifiers     string `json:"identifiers"`
		Name            string `json:"name"`
		SoftwareVersion string `json:"sw_version"`
		Model           string `json:"model"`
		Manufacturer    string `json:"manufacturer"`
	} `json:"device"`

	subscribersOnce sync.Once
	subscribers     []mqtt.Subscriber
	state           atomic.Value
}

func (c *Base) GetID() string {
	return c.ID
}

func (c *Base) GetType() ComponentType {
	return c.Type
}

func (c *Base) GetUniqueID() string {
	return c.UniqueID
}

func (c *Base) GetName() string {
	return c.Name
}

func (c *Base) GetState() string {
	if s := c.state.Load(); s != nil {
		return s.(string)
	}

	return ""
}

func (c *Base) GetCommandTopic() mqtt.Topic {
	return c.CommandTopic
}

func (c *Base) Subscribers() []mqtt.Subscriber {
	c.subscribersOnce.Do(func() {
		c.subscribers = make([]mqtt.Subscriber, 0)

		if c.StateTopic != "" {
			c.subscribers = append(c.subscribers, mqtt.NewSubscriber(c.StateTopic, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				c.state.Store(message.String())
				return nil
			}))
		}
	})

	return c.subscribers
}
