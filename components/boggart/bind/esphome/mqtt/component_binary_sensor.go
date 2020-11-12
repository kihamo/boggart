package mqtt

import (
	"bytes"
	"encoding/json"

	"github.com/kihamo/boggart/components/mqtt"
)

type ComponentBinarySensor struct {
	*ComponentBase

	data struct {
		DeviceClass string `json:"device_class"`
		PayloadOn   string `json:"payload_on"`
		PayloadOff  string `json:"payload_off"`
	}
}

func NewComponentBinarySensor(id string) *ComponentBinarySensor {
	component := &ComponentBinarySensor{
		ComponentBase: NewComponentBase(id, ComponentTypeBinarySensor),
	}
	component.setState = component.SetState

	return component
}

func (c *ComponentBinarySensor) SetState(message mqtt.Message) error {
	var val float64

	if c.ConvertState(message.String()) {
		val = 1
	}

	metricState.With("mac", c.Device().MAC().String()).With("component", c.ID()).Set(val)

	return c.ComponentBase.SetState(message)
}

func (c *ComponentBinarySensor) DeviceClass() string {
	return c.data.DeviceClass
}

func (c *ComponentBinarySensor) PayloadOn() string {
	return c.data.PayloadOn
}

func (c *ComponentBinarySensor) PayloadOff() string {
	return c.data.DeviceClass
}

func (c *ComponentBinarySensor) ConvertState(state string) bool {
	if c.data.PayloadOn != "" && state == c.data.PayloadOn {
		return true
	}

	if bytes.Equal([]byte(state), stateON) {
		return true
	}

	return false
}

func (c *ComponentBinarySensor) UnmarshalJSON(b []byte) error {
	if err := c.ComponentBase.UnmarshalJSON(b); err != nil {
		return err
	}

	return json.Unmarshal(b, &c.data)
}
