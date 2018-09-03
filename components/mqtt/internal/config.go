package internal

import (
	m "github.com/eclipse/paho.mqtt.golang"
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
		config.NewVariable(mqtt.ConfigUsername, config.ValueTypeString).
			WithUsage("Username").
			WithGroup("Connect").
			WithEditable(true),
		config.NewVariable(mqtt.ConfigPassword, config.ValueTypeString).
			WithUsage("Password").
			WithGroup("Connect").
			WithView([]string{config.ViewPassword}).
			WithEditable(true),
		config.NewVariable(mqtt.ConfigConnectionTimeout, config.ValueTypeDuration).
			WithUsage("Timeout").
			WithGroup("Connect").
			WithDefault(m.NewClientOptions().ConnectTimeout),
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
	c.initClient()
}
