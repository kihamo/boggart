package mqtt

import (
	"context"
	"encoding/json"
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
	data struct {
		Icon                string     `json:"icon"`
		Name                string     `json:"name"`
		StateTopic          mqtt.Topic `json:"state_topic"`
		CommandTopic        mqtt.Topic `json:"command_topic"`
		AvailabilityTopic   mqtt.Topic `json:"availability_topic"`
		PayloadAvailable    string     `json:"payload_available"`
		PayloadNotAvailable string     `json:"payload_not_available"`
		UniqueID            string     `json:"unique_id"`
		Device              Device     `json:"device"`
	}

	id  string
	typ ComponentType

	subscribersOnce sync.Once
	subscribers     []mqtt.Subscriber
	state           atomic.Value
	setState        func(mqtt.Message) error
}

func NewComponentBase(id string, t ComponentType) *ComponentBase {
	component := &ComponentBase{
		id:  id,
		typ: t,
	}
	component.setState = component.SetState

	return component
}

func (c *ComponentBase) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &c.data)
}

func (c *ComponentBase) ID() string {
	return c.id
}

func (c *ComponentBase) Type() ComponentType {
	return c.typ
}

func (c *ComponentBase) UniqueID() string {
	return c.data.UniqueID
}

func (c *ComponentBase) Name() string {
	return c.data.Name
}

func (c *ComponentBase) State() interface{} {
	if s := c.state.Load(); s != nil {
		return s.(string)
	}

	return ""
}

func (c *ComponentBase) StateTopic() mqtt.Topic {
	return c.data.StateTopic
}

func (c *ComponentBase) CommandTopic() mqtt.Topic {
	return c.data.CommandTopic
}

func (c *ComponentBase) AvailabilityTopic() mqtt.Topic {
	return c.data.AvailabilityTopic
}

func (c *ComponentBase) Device() Device {
	return c.data.Device
}

// nolint:interfacer
func (c *ComponentBase) SetState(message mqtt.Message) error {
	state := message.String()

	c.state.Store(state)

	if val, err := strconv.ParseFloat(state, 64); err == nil {
		metricState.With("mac", c.Device().MAC().String()).With("component", c.ID()).Set(val)
	}

	return nil
}

func (c *ComponentBase) CommandToPayload(cmd interface{}) interface{} {
	return cmd
}

func (c *ComponentBase) Subscribers() []mqtt.Subscriber {
	c.subscribersOnce.Do(func() {
		c.subscribers = make([]mqtt.Subscriber, 0)

		if topic := c.StateTopic(); topic != "" {
			c.subscribers = append(c.subscribers, mqtt.NewSubscriber(topic, 0, func(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
				return c.setState(message)
			}))
		}
	})

	return c.subscribers
}
