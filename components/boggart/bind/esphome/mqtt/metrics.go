package mqtt

import (
	"github.com/kihamo/snitch"
)

var (
	metricState              = snitch.NewGauge("state", "Component state")
	metricSensorValue        = snitch.NewGauge("sensor_value", "Export metric esphome_sensor_value")
	metricSensorFailed       = snitch.NewGauge("sensor_failed", "Export metric esphome_sensor_failed")
	metricBinarySensorValue  = snitch.NewGauge("binary_sensor_value", "Export metric esphome_binary_sensor_value")
	metricBinarySensorFailed = snitch.NewGauge("binary_sensor_failed", "Export metric esphome_binary_sensor_failed")
	metricFanValue           = snitch.NewGauge("fan_value", "Export metric esphome_fan_value")
	metricFanFailed          = snitch.NewGauge("fan_failed", "Export metric esphome_fan_failed")
	metricFanSpeed           = snitch.NewGauge("fan_speed", "Export metric esphome_fan_speed")
	metricFanOscillation     = snitch.NewGauge("fan_oscillation", "Export metric esphome_fan_oscillation")
	metricLightState         = snitch.NewGauge("light_state", "Export metric esphome_light_state")
	metricLightColor         = snitch.NewGauge("light_color", "Export metric esphome_light_color")
	metricLightEffectActive  = snitch.NewGauge("light_effect_active", "Export metric esphome_light_effect_active")
	metricCoverValue         = snitch.NewGauge("cover_value", "Export metric esphome_cover_value")
	metricCoverFailed        = snitch.NewGauge("cover_failed", "Export metric esphome_cover_failed")
	metricSwitchValue        = snitch.NewGauge("cover_value", "Export metric esphome_switch_value")
	metricSwitchFailed       = snitch.NewGauge("cover_failed", "Export metric esphome_switch_failed")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	mac := b.Meta().MACAsString()
	if mac == "" {
		return
	}

	metricState.With("mac", mac).Describe(ch)

	if b.IP() != nil {
		metricSensorValue.With("mac", mac).Describe(ch)
		metricSensorFailed.With("mac", mac).Describe(ch)
		metricBinarySensorValue.With("mac", mac).Describe(ch)
		metricBinarySensorFailed.With("mac", mac).Describe(ch)
		metricFanValue.With("mac", mac).Describe(ch)
		metricFanFailed.With("mac", mac).Describe(ch)
		metricFanSpeed.With("mac", mac).Describe(ch)
		metricFanOscillation.With("mac", mac).Describe(ch)
		metricLightState.With("mac", mac).Describe(ch)
		metricLightColor.With("mac", mac).Describe(ch)
		metricLightEffectActive.With("mac", mac).Describe(ch)
		metricCoverValue.With("mac", mac).Describe(ch)
		metricCoverFailed.With("mac", mac).Describe(ch)
		metricSwitchValue.With("mac", mac).Describe(ch)
		metricSwitchFailed.With("mac", mac).Describe(ch)
	}
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	mac := b.Meta().MACAsString()
	if mac == "" {
		return
	}

	metricState.With("mac", mac).Collect(ch)

	if b.IP() != nil {
		metricSensorValue.With("mac", mac).Collect(ch)
		metricSensorFailed.With("mac", mac).Collect(ch)
		metricBinarySensorValue.With("mac", mac).Collect(ch)
		metricBinarySensorFailed.With("mac", mac).Collect(ch)
		metricFanValue.With("mac", mac).Collect(ch)
		metricFanFailed.With("mac", mac).Collect(ch)
		metricFanSpeed.With("mac", mac).Collect(ch)
		metricFanOscillation.With("mac", mac).Collect(ch)
		metricLightState.With("mac", mac).Collect(ch)
		metricLightColor.With("mac", mac).Collect(ch)
		metricLightEffectActive.With("mac", mac).Collect(ch)
		metricCoverValue.With("mac", mac).Collect(ch)
		metricCoverFailed.With("mac", mac).Collect(ch)
		metricSwitchValue.With("mac", mac).Collect(ch)
		metricSwitchFailed.With("mac", mac).Collect(ch)
	}
}
