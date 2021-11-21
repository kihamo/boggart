package mqtt

import (
	"bytes"
	"encoding/json"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/mqtt"
)

var (
	stateON = []byte(`ON`)
)

type ComponentSwitch struct {
	*componentBase

	// https://github.com/esphome/esphome/blob/2021.11.1/esphome/components/mqtt/mqtt_switch.cpp#L47
	data struct {
		Optimistic bool `json:"optimistic"`
	}

	state *atomic.BoolNull
}

func NewComponentSwitch(id string, message mqtt.Message) *ComponentSwitch {
	return &ComponentSwitch{
		componentBase: newComponentBase(id, ComponentTypeSwitch, message),
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
		metricState.With("mac", c.DeviceInfo().MAC().String()).With("component", c.ID()).Set(1)

		return nil
	}

	c.state.False()
	metricState.With("mac", c.DeviceInfo().MAC().String()).With("component", c.ID()).Set(0)

	return nil
}

func (c *ComponentSwitch) UnmarshalJSON(b []byte) error {
	if err := c.componentBase.UnmarshalJSON(b); err != nil {
		return err
	}

	return json.Unmarshal(b, &c.data)
}
