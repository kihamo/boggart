package internal

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			boggart.ConfigPulsarEnabled,
			config.ValueTypeBool,
			true,
			"Enabled pulsar provider",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigPulsarRepeatInterval,
			config.ValueTypeDuration,
			time.Minute*15,
			"Repeat interval for pulsar provider",
			false,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigPulsarSerialAddress,
			config.ValueTypeString,
			pulsar.DefaultSerialAddress,
			"Serial port address",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigPulsarSerialTimeout,
			config.ValueTypeDuration,
			pulsar.DefaultTimeout,
			"Serial port address",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigPulsarDeviceAddress,
			config.ValueTypeString,
			nil,
			"Pulsar address HEX value (AABBCCDD). If empty system try to find device",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigSoftVideoEnabled,
			config.ValueTypeBool,
			true,
			"Enabled SoftVideo provider",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigSoftVideoRepeatInterval,
			config.ValueTypeDuration,
			time.Hour*12,
			"Repeat interval for SoftVideo provider",
			false,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigSoftVideoLogin,
			config.ValueTypeString,
			nil,
			"Login of SoftVideo account",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigSoftVideoPassword,
			config.ValueTypeString,
			nil,
			"Password of SoftVideo account",
			true,
			[]string{config.ViewPassword},
			nil),
	}
}
