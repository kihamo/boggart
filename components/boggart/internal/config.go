package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(boggart.ConfigConfigYAML, config.ValueTypeString).
			WithUsage("Absolute path to YAML config").
			WithGroup("Config").
			WithDefault("config.yaml"),
	}
}
