package internal

import (
	"strconv"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(mqtt.ConfigServers, config.ValueTypeString).
			WithUsage("Server").
			WithGroup("Connect").
			WithDefault("tcp://localhost:1883").
			WithEditable(true),
		config.NewVariable(mqtt.ConfigClientID, config.ValueTypeString).
			WithUsage("Topic").
			WithGroup("Client ID").
			WithDefaultFunc(func() interface{} {
				name := mqtt.NameReplace(c.application.Name())
				if len(name) > 10 {
					name = name[0:9]
				}

				return name + "_v" + strconv.FormatInt(c.application.BuildDate().Unix(), 10)
			}),
		config.NewVariable(mqtt.ConfigUsername, config.ValueTypeString).
			WithUsage("Username").
			WithGroup("Connect").
			WithEditable(true),
		config.NewVariable(mqtt.ConfigPassword, config.ValueTypeString).
			WithUsage("Password").
			WithGroup("Connect").
			WithView([]string{config.ViewPassword}).
			WithEditable(true),
		config.NewVariable(mqtt.ConfigConnectionAttempts, config.ValueTypeInt64).
			WithUsage("Attempts (if 0 unlimited)").
			WithGroup("Connect").
			WithDefault(mqtt.DefaultConnectionAttempts),
		config.NewVariable(mqtt.ConfigConnectionTimeout, config.ValueTypeDuration).
			WithUsage("Timeout").
			WithGroup("Connect").
			WithDefault(m.NewClientOptions().ConnectTimeout),
		config.NewVariable(mqtt.ConfigClearSession, config.ValueTypeBool).
			WithUsage("Clear session").
			WithGroup("Connect").
			WithDefault(true),
		config.NewVariable(mqtt.ConfigLWTEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Last Will and Testament").
			WithDefault(false),
		config.NewVariable(mqtt.ConfigLWTTopic, config.ValueTypeString).
			WithUsage("Topic").
			WithGroup("Last Will and Testament").
			WithDefaultFunc(func() interface{} {
				return boggart.ComponentName + "/" + mqtt.NameReplace(c.application.Name()) + "/status"
			}),
		config.NewVariable(mqtt.ConfigLWTPayload, config.ValueTypeString).
			WithUsage("Payload").
			WithGroup("Last Will and Testament").
			WithDefault("0"),
		config.NewVariable(mqtt.ConfigLWTQOS, config.ValueTypeInt64).
			WithUsage("OQS").
			WithGroup("Last Will and Testament").
			WithDefault(0).
			WithView([]string{config.ViewEnum}).
			WithViewOptions(map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{"0", "At most once"},
					{"1", "At least once"},
					{"2", "Exactly once"},
				},
			}),
		config.NewVariable(mqtt.ConfigLWTRetained, config.ValueTypeBool).
			WithUsage("Retained").
			WithGroup("Last Will and Testament").
			WithDefault(false),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{
			mqtt.ConfigServers,
			mqtt.ConfigUsername,
			mqtt.ConfigPassword,
		}, c.watchConnect),
	}
}

func (c *Component) watchConnect(_ string, _ interface{}, _ interface{}) {
	if err := c.initClient(); err == nil {
		c.initSubscribers()
	}
}
