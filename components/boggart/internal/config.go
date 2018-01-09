package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			boggart.ConfigPulsarSerialPath,
			config.ValueTypeString,
			"/dev/ttyUSB0",
			"Pulsar device path",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigPulsarAddress,
			config.ValueTypeString,
			nil,
			"Pulsar address HEX value",
			true,
			nil,
			nil),
	}
}
