package mqtt

import (
	"context"
	"net"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
)

type Device struct {
	Identifiers     string `json:"identifiers"`
	Name            string `json:"name"`
	SoftwareVersion string `json:"sw_version"`
	Model           string `json:"model"`
	Manufacturer    string `json:"manufacturer"`
}

func (d Device) MAC() (mac net.HardwareAddr) {
	id := d.Identifiers

	for i := 2; i < len(id); i += 3 {
		id = id[:i] + "-" + id[i:]
	}

	mac, _ = net.ParseMAC(id)

	return mac
}

type ComponentBase struct {
	ID                  string        `json:"-"`
	Type                ComponentType `json:"-"`
	Icon                string        `json:"icon"`
	Name                string        `json:"name"`
	StateTopic          mqtt.Topic    `json:"state_topic"`
	CommandTopic        mqtt.Topic    `json:"command_topic"`
	PayloadAvailable    string        `json:"payload_available"`
	PayloadNotAvailable string        `json:"payload_not_available"`
	UniqueID            string        `json:"unique_id"`
	Device              Device        `json:"device"`

	subscribersOnce sync.Once
	subscribers     []mqtt.Subscriber
	state           atomic.Value
	setState        func(mqtt.Message) error
}

func NewComponentBase(id string, t ComponentType) *ComponentBase {
	component := &ComponentBase{
		ID:   id,
		Type: t,
	}
	component.setState = component.SetState

	return component
}

func (c *ComponentBase) GetID() string {
	return c.ID
}

func (c *ComponentBase) GetType() ComponentType {
	return c.Type
}

func (c *ComponentBase) GetUniqueID() string {
	return c.UniqueID
}

func (c *ComponentBase) GetName() string {
	return c.Name
}

func (c *ComponentBase) GetState() interface{} {
	if s := c.state.Load(); s != nil {
		return s.(string)
	}

	return ""
}

func (c *ComponentBase) GetCommandTopic() mqtt.Topic {
	return c.CommandTopic
}

func (c *ComponentBase) GetDevice() Device {
	return c.Device
}

func (c *ComponentBase) SetState(message mqtt.Message) error {
	c.state.Store(message.String())

	if val, err := strconv.ParseFloat(message.String(), 64); err == nil {
		metricState.With("serial_number", c.Device.MAC().String()).With("component", c.ID).Set(val)
	}

	return nil
}

func (c *ComponentBase) CommandToPayload(cmd interface{}) interface{} {
	return cmd
}

func (c *ComponentBase) Subscribers() []mqtt.Subscriber {
	c.subscribersOnce.Do(func() {
		c.subscribers = make([]mqtt.Subscriber, 0)

		if c.StateTopic != "" {
			c.subscribers = append(c.subscribers, mqtt.NewSubscriber(c.StateTopic, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return c.setState(message)
			}))
		}
	})

	return c.subscribers
}
