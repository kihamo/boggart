package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
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
		config.NewVariable(boggart.ConfigBoggartBindID, config.ValueTypeString).
			WithUsage("Boggart bind ID").
			WithDefaultFunc(func() interface{} {
				return mqtt.NameReplace(c.application.Name())
			}),
		config.NewVariable(boggart.ConfigMQTTTopicBindStatus, config.ValueTypeString).
			WithUsage("Boggart MQTT topic bind status").
			WithDefaultFunc(func() interface{} {
				return boggart.ComponentName + "/bind/+/status"
			}).
			WithEditable(true),
	}
}
