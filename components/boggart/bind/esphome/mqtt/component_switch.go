package mqtt

import (
	"bytes"

	"github.com/kihamo/boggart/components/mqtt"
)

var (
	stateON = []byte(`ON`)
)

type ComponentSwitch struct {
	*ComponentBase
}

func NewComponentSwitch(id string) *ComponentSwitch {
	component := &ComponentSwitch{
		ComponentBase: NewComponentBase(id, ComponentTypeSwitch),
	}
	component.setState = component.SetState

	return component
}

func (c *ComponentSwitch) SetState(message mqtt.Message) {
	c.ComponentBase.SetState(message)

	var val float64
	if bytes.Equal(message.Payload(), stateON) {
		val = 1
	}

	metricState.With("serial_number", c.Device.MAC().String()).With("component", c.ID).Set(val)
}
