package mqtt

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{
		ip: atomic.NewValue(),
	}
}

func (t Type) DashboardTemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		"is_allow_light_state_field": IsAllowLightStateField,
	}
}

func IsAllowLightStateField(m ComponentLightColorModes, field string) bool {
	switch field {
	case "color-modes":
		return !m.IsOnOff()
	case "brightness":
		return m.IsBrightness() || m.IsWhite() || m.IsColorTemperature() || m.IsColdWarmWhite() ||
			m.IsRGB() || m.IsRGBWhite() || m.IsRGBColorTemperature() || m.IsRGBColdWarmWhite()
	case "red", "green", "blue":
		return m.IsRGB() || m.IsRGBWhite() || m.IsRGBColorTemperature() || m.IsRGBColdWarmWhite()
	case "white", "warm-white":
		return m.IsWhite() || m.IsColdWarmWhite() || m.IsRGBWhite() || m.IsColdWarmWhite() || m.IsRGBColdWarmWhite()
	case "cold-white":
		return m.IsColdWarmWhite() || m.IsRGBColdWarmWhite()
	case "color-temperature":
		return m.IsColorTemperature() || m.IsColdWarmWhite() || m.IsRGBColorTemperature() || m.IsRGBColdWarmWhite()
	}

	return true
}
