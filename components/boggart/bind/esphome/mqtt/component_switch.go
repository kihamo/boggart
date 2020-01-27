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

func (c *ComponentSwitch) SetState(message mqtt.Message) error {
	var val float64
	if bytes.Equal(message.Payload(), stateON) {
		val = 1
	}

	metricState.With("mac", c.Device.MAC().String()).With("component", c.ID).Set(val)

	return c.ComponentBase.SetState(message)
}
