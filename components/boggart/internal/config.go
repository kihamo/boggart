package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
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
	}
}
