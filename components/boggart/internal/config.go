package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(boggart.ConfigConfigYAML, config.ValueTypeString).
			WithUsage("Absolute path to YAML config").
			WithDefault("config.yaml"),
		config.NewVariable(boggart.ConfigAccessKeys, config.ValueTypeString).
			WithUsage("Access keys").
			WithEditable(true).
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a key"}),
	}
}
