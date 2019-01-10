package internal

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(boggart.ConfigListenerTelegramChats, config.ValueTypeString).
			WithUsage("Chats for messages").
			WithGroup("Listener Telegram").
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a chat ID"}),
		config.NewVariable(boggart.ConfigRS485Address, config.ValueTypeString).
			WithUsage("Serial port address").
			WithGroup("RS485 protocol").
			WithDefault(rs485.DefaultSerialAddress),
		config.NewVariable(boggart.ConfigRS485Timeout, config.ValueTypeDuration).
			WithUsage("Serial port timeout").
			WithGroup("RS485 protocol").
			WithDefault(rs485.DefaultTimeout),
		config.NewVariable(boggart.ConfigMercuryRepeatInterval, config.ValueTypeDuration).
			WithUsage("Repeat interval").
			WithGroup("Mercury devices").
			WithDefault(time.Minute * 2),
		config.NewVariable(boggart.ConfigMercuryDeviceAddress, config.ValueTypeString).
			WithUsage("Device address in format XXXXXX (last 6 digits of device serial number)").
			WithGroup("Mercury devices"),
		config.NewVariable(boggart.ConfigMQTTOwnTracksEnabled, config.ValueTypeBool).
			WithUsage("OwnTracks enabled").
			WithGroup("MQTT subscribers"),
		config.NewVariable(boggart.ConfigMQTTWOLEnabled, config.ValueTypeBool).
			WithUsage("Wake-on-LAN enabled").
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
