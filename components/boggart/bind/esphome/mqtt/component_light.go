package mqtt

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
)

type ComponentLightState struct {
	Effect     string `json:"effect,omitempty"`
	State      string `json:"state,omitempty"`
	Brightness uint64 `json:"brightness,omitempty"`
	Color      struct {
		Red   uint64 `json:"r,omitempty"`
		Green uint64 `json:"g,omitempty"`
		Blue  uint64 `json:"b,omitempty"`
	} `json:"color,omitempty"`
	White            uint64 `json:"white_value,omitempty"`
	ColorTemperature uint64 `json:"color_temp,omitempty"`
	Flash            uint64 `json:"flash,omitempty"`
	Transition       uint64 `json:"transition,omitempty"`
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

func NewComponentLight(id string) *ComponentLight {
	return &ComponentLight{
		componentBase: newComponentBase(id, ComponentTypeLight),
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

func (c *ComponentLight) UnmarshalJSON(b []byte) error {
	if err := c.componentBase.UnmarshalJSON(b); err != nil {
		return err
	}

	return json.Unmarshal(b, &c.data)
}
