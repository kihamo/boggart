package mqtt

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
)

type ComponentLightState struct {
	Effect     string `json:"effect"`
	State      string `json:"state"`
	Brightness uint64 `json:"brightness"`
	Color      struct {
		Red   uint64 `json:"r"`
		Green uint64 `json:"g"`
		Blue  uint64 `json:"b"`
	} `json:"color"`
	White            uint64 `json:"white_value"`
	ColorTemperature uint64 `json:"color_temp"`
	Flash            uint64 `json:"flash"`
	Transition       uint64 `json:"transition"`
}

func (s *ComponentLightState) String() string {
	return s.State
}

func (s *ComponentLightState) SetState(state bool) {
	if state {
		s.State = "ON"
	} else {
		s.State = "OFF"
	}
}

type ComponentLight struct {
	*componentBase

	data struct {
		Schema           string   `json:"schema"`
		Brightness       bool     `json:"brightness"`
		RGB              bool     `json:"rgb"`
		ColorTemperature bool     `json:"color_temp"`
		White            bool     `json:"white_value"`
		Effect           bool     `json:"effect"`
		EffectList       []string `json:"effect_list"`
	}

	state atomic.Value
}

func NewComponentLight(id string, discoveryTopic mqtt.Topic) *ComponentLight {
	return &ComponentLight{
		componentBase: newComponentBase(id, ComponentTypeLight, discoveryTopic),
	}
}

func (c *ComponentLight) State() interface{} {
	if s := c.state.Load(); s != nil {
		return s.(*ComponentLightState)
	}

	return nil
}

func (c *ComponentLight) StateFormat() string {
	if s := c.state.Load(); s != nil {
		return s.(*ComponentLightState).String()
	}

	return ""
}

func (c *ComponentLight) SetState(message mqtt.Message) error {
	var state ComponentLightState

	if err := message.JSONUnmarshal(&state); err != nil {
		return err
	}

	c.state.Store(&state)

	return nil
}

func (c *ComponentLight) CommandToPayload(cmd interface{}) interface{} {
	var state ComponentLightState

	if st, ok := cmd.(*ComponentLightState); ok {
		state = *st
	} else {
		if s := c.state.Load(); s != nil {
			state = *(s.(*ComponentLightState))
		}

		if st, ok := cmd.(bool); ok {
			state.SetState(st)
		} else {
			state.State = fmt.Sprintf("%v", cmd)
		}
	}

	payload, _ := json.Marshal(state)

	return payload
}

func (c *ComponentLight) Schema() string {
	return c.data.Schema
}

func (c *ComponentLight) Brightness() bool {
	return c.data.Brightness
}

func (c *ComponentLight) RGB() bool {
	return c.data.RGB
}

func (c *ComponentLight) ColorTemperature() bool {
	return c.data.ColorTemperature
}

func (c *ComponentLight) White() bool {
	return c.data.White
}

func (c *ComponentLight) Effect() bool {
	return c.data.Effect
}

func (c *ComponentLight) EffectList() []string {
	return c.data.EffectList
}

func (c *ComponentLight) UnmarshalJSON(b []byte) error {
	if err := c.componentBase.UnmarshalJSON(b); err != nil {
		return err
	}

	return json.Unmarshal(b, &c.data)
}
