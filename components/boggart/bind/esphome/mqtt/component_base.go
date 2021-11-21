package mqtt

import (
	"encoding/json"
	"net"

	"github.com/kihamo/boggart/components/mqtt"
)

type DeviceInfo struct {
	Identifiers     string `json:"identifiers"`
	Name            string `json:"name"`
	SoftwareVersion string `json:"sw_version"`
	Model           string `json:"model"`
	Manufacturer    string `json:"manufacturer"`
}

func (d DeviceInfo) MAC() (mac net.HardwareAddr) {
	id := d.Identifiers

	for i := 2; i < len(id); i += 3 {
		id = id[:i] + "-" + id[i:]
	}

	mac, _ = net.ParseMAC(id)

	return mac
}

// https://github.com/esphome/esphome/blob/2021.11.1/esphome/components/mqtt/mqtt_component.cpp#L54
type componentBase struct {
	data struct {
		Name                string     `json:"name"`
		EnabledByDefault    bool       `json:"enabled_by_default"`
		Icon                string     `json:"icon"`
		EntityCategory      string     `json:"entity_category"`
		StateTopic          mqtt.Topic `json:"state_topic"`
		CommandTopic        mqtt.Topic `json:"command_topic"`
		AvailabilityTopic   mqtt.Topic `json:"availability_topic"`
		PayloadAvailable    string     `json:"payload_available"`
		PayloadNotAvailable string     `json:"payload_not_available"`
		UniqueID            string     `json:"unique_id"`
		DeviceInfo          DeviceInfo `json:"device"`
	}

	id            string
	typ           ComponentType
	configMessage mqtt.Message
	subscribers   []mqtt.Subscriber
}

func newComponentBase(id string, t ComponentType, message mqtt.Message) *componentBase {
	return &componentBase{
		id:            id,
		typ:           t,
		configMessage: message,
		subscribers:   make([]mqtt.Subscriber, 0, 3),
	}
}

func (c *componentBase) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &c.data)
}

func (c *componentBase) ID() string {
	return c.id
}

func (c *componentBase) Type() ComponentType {
	return c.typ
}

func (c *componentBase) UniqueID() string {
	return c.data.UniqueID
}

func (c *componentBase) Name() string {
	return c.data.Name
}

func (c *componentBase) EntityCategory() string {
	return c.data.EntityCategory
}

func (c *componentBase) Icon() string {
	return c.data.Icon
}

func (c *componentBase) ConfigMessage() mqtt.Message {
	return c.configMessage
}

func (c *componentBase) StateTopic() mqtt.Topic {
	return c.data.StateTopic
}

func (c *componentBase) CommandTopic() mqtt.Topic {
	return c.data.CommandTopic
}

func (c *componentBase) AvailabilityTopic() mqtt.Topic {
	return c.data.AvailabilityTopic
}

func (c *componentBase) DeviceInfo() DeviceInfo {
	return c.data.DeviceInfo
}

func (c *componentBase) CommandToPayload(cmd interface{}) interface{} {
	return cmd
}

func (c *componentBase) Subscribe(subscribers ...mqtt.Subscriber) {
	c.subscribers = append(c.subscribers, subscribers...)
}

func (c *componentBase) Subscribers() []mqtt.Subscriber {
	return c.subscribers
}
