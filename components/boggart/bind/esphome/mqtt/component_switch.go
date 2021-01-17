package mqtt

import (
	"bytes"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/mqtt"
)

var (
	stateON = []byte(`ON`)
)

type ComponentSwitch struct {
	*componentBase

	state *atomic.BoolNull
}

func NewComponentSwitch(id string) *ComponentSwitch {
	return &ComponentSwitch{
		componentBase: newComponentBase(id, ComponentTypeSwitch),
		state:         atomic.NewBoolNull(),
	}
}

func (c *ComponentSwitch) State() interface{} {
	if c.state.IsNil() {
		return nil
	}

	return c.state.Load()
}

func (c *ComponentSwitch) StateFormat() string {
	if c.state.IsNil() {
		return ""
	}

	if c.state.IsTrue() {
		return "ON"
	}

	return "OFF"
}

func (c *ComponentSwitch) SetState(message mqtt.Message) error {
	payload := message.Payload()

	if bytes.Equal(payload, stateON) {
		c.state.True()
		metricState.With("mac", c.Device().MAC().String()).With("component", c.ID()).Set(1)

		return nil
	}

	c.state.False()
	metricState.With("mac", c.Device().MAC().String()).With("component", c.ID()).Set(0)

	return nil
}
