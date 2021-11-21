package mqtt

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/kihamo/boggart/components/mqtt"
)

type ComponentLightState struct {
	Effect    string `json:"effect"`
	State     string `json:"state"`
	ColorMode string `json:"color_mode"`
	Color     struct {
		Red   uint64 `json:"r"`
		Green uint64 `json:"g"`
		Blue  uint64 `json:"b"`
		White uint64 `json:"w"`
		Cold  uint64 `json:"c"`
	} `json:"color"`
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

type ComponentLightColorModes []string

func (m ComponentLightColorModes) IsMode(mode string) bool {
	for _, m := range []string(m) {
		if mode == m {
			return true
		}
	}

	return false
}

func (m ComponentLightColorModes) IsOnOff() bool {
	return m.IsMode("onoff")
}

func (m ComponentLightColorModes) IsBrightness() bool {
	return m.IsMode("brightness")
}

func (m ComponentLightColorModes) IsColorTemperature() bool {
	return m.IsMode("color_temp")
}

func (m ComponentLightColorModes) IsColdWarmWhite() bool {
	return m.IsMode("color_temp")
}

func (m ComponentLightColorModes) IsRGB() bool {
	return m.IsMode("rgb")
}

func (m ComponentLightColorModes) IsRGBWhite() bool {
	return m.IsMode("rgbw")
}

func (m ComponentLightColorModes) IsRGBColorTemperature() bool {
	return m.IsMode("rgbw")
}

func (m ComponentLightColorModes) IsRGBColdWarmWhite() bool {
	return m.IsMode("rgbww")
}

type ComponentLight struct {
	*componentBase

	// https://github.com/esphome/esphome/blob/2021.11.1/esphome/components/mqtt/mqtt_light.cpp#L39
	data struct {
		Schema     string                   `json:"schema"`
		ColorModes ComponentLightColorModes `json:"supported_color_modes"`
		Effect     bool                     `json:"effect"`
		EffectList []string                 `json:"effect_list"`
	}

	state atomic.Value
}

func NewComponentLight(id string, message mqtt.Message) *ComponentLight {
	return &ComponentLight{
		componentBase: newComponentBase(id, ComponentTypeLight, message),
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

func (c *ComponentLight) ColorModes() ComponentLightColorModes {
	return c.data.ColorModes
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
