package internal

import (
	"github.com/kihamo/boggart/components/syslog"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(syslog.ConfigHost, config.ValueTypeString).
			WithUsage("Host").
			WithGroup("Listen").
			WithDefault("localhost"),
		config.NewVariable(syslog.ConfigPort, config.ValueTypeInt).
			WithUsage("Port number").
			WithGroup("Listen").
			WithDefault(514),
	}
}
