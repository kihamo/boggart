package mqtt

import (
	"bytes"

	"github.com/kihamo/boggart/components/mqtt"
)

type ComponentBinarySensor struct {
	*ComponentBase

	DeviceClass string `json:"device_class"`
	PayloadOn   string `json:"payload_on"`
	PayloadOff  string `json:"payload_off"`
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

	if c.PayloadOn != "" && message.String() == c.PayloadOn {
		val = 1
	} else if bytes.Equal(message.Payload(), stateON) {
		val = 1
	}

	metricState.With("mac", c.Device().MAC().String()).With("component", c.ID()).Set(val)

	return c.ComponentBase.SetState(message)
}
