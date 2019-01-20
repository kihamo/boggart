package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(boggart.ConfigListenerTelegramChats, config.ValueTypeString).
			WithUsage("Chats for messages").
			WithGroup("Listener Telegram").
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a chat ID"}),
		config.NewVariable(boggart.ConfigMQTTOwnTracksEnabled, config.ValueTypeBool).
			WithUsage("OwnTracks enabled").
			WithGroup("MQTT subscribers"),
		config.NewVariable(boggart.ConfigMQTTAnnotationsEnabled, config.ValueTypeBool).
			WithUsage("Annotations enabled").
			WithGroup("MQTT subscribers"),
		config.NewVariable(boggart.ConfigMQTTMessengersEnabled, config.ValueTypeBool).
			WithUsage("Messengers enabled").
			WithGroup("MQTT subscribers"),
		config.NewVariable(boggart.ConfigConfigYAML, config.ValueTypeString).
			WithUsage("Absolute path to YAML config").
			WithGroup("Config").
			WithDefault("config.yaml"),
	}
}
