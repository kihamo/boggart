package internal

import (
	"github.com/kihamo/boggart/components/syslog"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			syslog.ConfigHost,
			config.ValueTypeString,
			"localhost",
			"Host",
			false,
			"Listen",
			nil,
			nil),
		config.NewVariable(
			syslog.ConfigPort,
			config.ValueTypeInt,
			514,
			"Port number",
			false,
			"Listen",
			nil,
			nil),
	}
}
