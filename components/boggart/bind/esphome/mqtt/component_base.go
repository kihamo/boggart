package mqtt

import (
	"encoding/json"
	"net"

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

type componentBase struct {
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
}

func newComponentBase(id string, t ComponentType) *componentBase {
	return &componentBase{
		id:  id,
		typ: t,
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

func (c *componentBase) Icon() string {
	return c.data.Icon
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

func (c *componentBase) Device() Device {
	return c.data.Device
}

func (c *componentBase) CommandToPayload(cmd interface{}) interface{} {
	return cmd
}
