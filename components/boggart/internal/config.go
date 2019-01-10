package internal

import (
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
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
		config.NewVariable(boggart.ConfigPulsarEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Pulsar devices"),
		config.NewVariable(boggart.ConfigPulsarRepeatInterval, config.ValueTypeDuration).
			WithUsage("Repeat interval").
			WithGroup("Pulsar devices").
			WithDefault(time.Minute * 15),
		config.NewVariable(boggart.ConfigPulsarHeatMeterAddress, config.ValueTypeString).
			WithUsage("Device address HEX value (AABBCCDD). If empty system try to find device").
			WithGroup("Pulsar devices"),
		config.NewVariable(boggart.ConfigPulsarColdWaterSerialNumber, config.ValueTypeString).
			WithUsage("Device address of cold water meter in format AABBCCDD").
			WithGroup("Pulsar devices"),
		config.NewVariable(boggart.ConfigPulsarColdWaterPulseInput, config.ValueTypeUint64).
			WithUsage("Input number of cold water meter in heat meter").
			WithGroup("Pulsar devices").
			WithDefault(pulsar.Input1).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{strconv.FormatUint(pulsar.Input1, 10), "#1"},
					{strconv.FormatUint(pulsar.Input2, 10), "#2"},
				},
			}),
		config.NewVariable(boggart.ConfigPulsarColdWaterStartValue, config.ValueTypeFloat64).
			WithUsage("Start value of cold water meter (in m3)").
			WithGroup("Pulsar devices").
			WithDefault(0),
		config.NewVariable(boggart.ConfigPulsarHotWaterSerialNumber, config.ValueTypeString).
			WithUsage("Device address of hot water meter in format AABBCCDD").
			WithGroup("Pulsar devices"),
		config.NewVariable(boggart.ConfigPulsarHotWaterPulseInput, config.ValueTypeUint64).
			WithUsage("Input number of hot water meter").
			WithGroup("Pulsar devices").
			WithDefault(pulsar.Input2).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{strconv.FormatUint(pulsar.Input1, 10), "#1"},
					{strconv.FormatUint(pulsar.Input2, 10), "#2"},
				},
			}),
		config.NewVariable(boggart.ConfigPulsarHotWaterStartValue, config.ValueTypeFloat64).
			WithUsage("Start value of hot water (in m3)").
			WithGroup("Pulsar devices").
			WithDefault(0),
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
