package mqtt

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
)

const (
	SchemaJSON = "json"
)

type ComponentLightState struct {
	State      string  `json:"state"`
	ColorMode  *string `json:"color_mode,omitempty"` // read only
	Brightness *uint64 `json:"brightness,omitempty"`
	Color      struct {
		Red       *uint64 `json:"r,omitempty"`
		Green     *uint64 `json:"g,omitempty"`
		Blue      *uint64 `json:"b,omitempty"`
		White     *uint64 `json:"w,omitempty"`
		ColdWhite *uint64 `json:"c,omitempty"`
	} `json:"color"`
	White            *uint64 `json:"white_value,omitempty"` // Deprecated: legacy API, use Color.White
	ColorTemperature *uint64 `json:"color_temp,omitempty"`
	Flash            *uint64 `json:"flash,omitempty"`      // write only
	Transition       *uint64 `json:"transition,omitempty"` // write only
	Effect           *string `json:"effect,omitempty"`
}

/*
brightness         : BRIGHTNESS, WHITE, COLOR_TEMPERATURE, COLD_WARM_WHITE, RGB, RGB_WHITE, RGB_COLOR_TEMPERATURE, RGB_COLD_WARM_WHITE
red                : RGB, RGB_WHITE, RGB_COLOR_TEMPERATURE, RGB_COLD_WARM_WHITE
green              : RGB, RGB_WHITE, RGB_COLOR_TEMPERATURE, RGB_COLD_WARM_WHITE
blue               : RGB, RGB_WHITE, RGB_COLOR_TEMPERATURE, RGB_COLD_WARM_WHITE
white      (w)     : WHITE, COLD_WARM_WHITE, RGB_WHITE
cold_white (c)     : COLD_WARM_WHITE, RGB_COLD_WARM_WHITE
warm_white (w)     : COLD_WARM_WHITE, RGB_COLD_WARM_WHITE
color_temperature  : COLOR_TEMPERATURE, COLD_WARM_WHITE, RGB_COLOR_TEMPERATURE, RGB_COLD_WARM_WHITE
color_brightness   : RGB, RGB_WHITE, RGB_COLOR_TEMPERATURE, RGB_COLD_WARM_WHITE
*/

func (s *ComponentLightState) String() string {
	if s.State != "OFF" && s.ColorMode != nil {
		return s.State + " " + *s.ColorMode
	}

	return s.State
}

func (s *ComponentLightState) SetState(state bool) {
	if state {
		s.State = "ON"
	} else {
		s.State = "OFF"
	}
}

func (s *ComponentLightState) SetColorMode(value string) {
	s.ColorMode = &value
}

func (s *ComponentLightState) SetBrightness(value uint64) {
	s.Brightness = &value
}

func (s *ComponentLightState) SetRed(value uint64) {
	s.Color.Red = &value
}

func (s *ComponentLightState) SetGreen(value uint64) {
	s.Color.Green = &value
}

func (s *ComponentLightState) SetBlue(value uint64) {
	s.Color.Blue = &value
}

func (s *ComponentLightState) SetWhite(value uint64) {
	s.Color.White = &value
}

func (s *ComponentLightState) SetColdWhite(value uint64) {
	s.Color.ColdWhite = &value
}

func (s *ComponentLightState) SetColorTemperature(value uint64) {
	s.ColorTemperature = &value
}

func (s *ComponentLightState) SetEffect(value string) {
	s.Effect = &value
}

func (s *ComponentLightState) SetFlash(value time.Duration) {
	if value > 0 {
		s.Flash = &[]uint64{uint64(value.Seconds())}[0]
	}
}

func (s *ComponentLightState) SetTransition(value time.Duration) {
	if value > 0 {
		s.Transition = &[]uint64{uint64(value.Seconds())}[0]
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

func (m ComponentLightColorModes) IsWhite() bool {
	return m.IsMode("white")
}

func (m ComponentLightColorModes) IsColorTemperature() bool {
	return m.IsMode("color_temp")
}

func (m ComponentLightColorModes) IsColdWarmWhite() bool {
	return m.IsMode("cwww")
}

func (m ComponentLightColorModes) IsRGB() bool {
	return m.IsMode("rgb")
}

func (m ComponentLightColorModes) IsRGBWhite() bool {
	return m.IsMode("rgbw")
}

func (m ComponentLightColorModes) IsRGBColorTemperature() bool {
	return m.IsMode("rgbct")
}

func (m ComponentLightColorModes) IsRGBColdWarmWhite() bool {
	return m.IsMode("rgbww")
}

type ComponentLight struct {
	*componentBase

	// https://github.com/esphome/esphome/blob/2021.11.1/esphome/components/mqtt/mqtt_light.cpp#L39
	data struct {
		Schema     string                   `json:"schema"`
		ColorMode  bool                     `json:"color_mode"`
		ColorModes ComponentLightColorModes `json:"supported_color_modes"`
		Brightness bool                     `json:"brightness"` // Deprecated: legacy API
		Effect     bool                     `json:"effect"`
		EffectList []string                 `json:"effect_list"`
	}

	state    atomic.Value
	stateRaw atomic.Value
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

func (c *ComponentLight) StateRaw() string {
	if s := c.stateRaw.Load(); s != nil {
		return s.(string)
	}

	return ""
}

func (c *ComponentLight) SetState(message mqtt.Message) error {
	var state ComponentLightState

	if err := message.JSONUnmarshal(&state); err != nil {
		return err
	}

	c.state.Store(&state)
	c.stateRaw.Store(message.String())

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

func (c *ComponentLight) ColorMode() bool {
	return c.data.ColorMode
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
