package internal

import (
	"strconv"
	"time"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	opts := m.NewClientOptions()

	return []config.Variable{
		config.NewVariable(mqtt.ConfigClientID, config.ValueTypeString).
			WithUsage("Client ID to be used by this client when connecting to the MQTT broker. According to the MQTT v3.1 specification, a client id mus be no longer than 23 characters.").
			WithGroup("Credentials").
			WithDefaultFunc(func() interface{} {
				name := mqtt.NameReplace(c.application.Name())
				if len(name) > 10 {
					name = name[0:9]
				}

				return name + "_v" + strconv.FormatInt(c.application.BuildDate().Unix(), 10)
			}),
		config.NewVariable(mqtt.ConfigUsername, config.ValueTypeString).
			WithUsage("Username").
			WithGroup("Credentials").
			WithEditable(true),
		config.NewVariable(mqtt.ConfigPassword, config.ValueTypeString).
			WithUsage("Password").
			WithGroup("Credentials").
			WithView([]string{config.ViewPassword}).
			WithEditable(true),
		config.NewVariable(mqtt.ConfigClearSession, config.ValueTypeBool).
			WithUsage("Flag in the connect message when this client connects to an MQTT broker. By setting this flag, you are indicating that no messages saved by the broker for this client should be delivered. Any messages that were going to be sent by this client before diconnecting previously but didn't will not be sent upon connecting to the broker.").
			WithGroup("Others").
			WithDefault(opts.CleanSession),
		config.NewVariable(mqtt.ConfigResumeSubs, config.ValueTypeBool).
			WithUsage("Will enable resuming of stored (un)subscribe messages when connecting but not reconnecting if CleanSession is false. Otherwise these messages are discarded.").
			WithGroup("Others").
			WithDefault(opts.ResumeSubs),
		config.NewVariable(mqtt.ConfigOrderMatters, config.ValueTypeBool).
			WithUsage("Will set the message routing to guarantee order within each QoS level. By default, this value is true. If set to false, this flag indicates that messages can be delivered asynchronously from the client to the application and possibly arrive out of order. Specifically, the message handler is called in its own go routine..").
			WithGroup("Others").
			WithDefault(opts.ResumeSubs),
		config.NewVariable(mqtt.ConfigServers, config.ValueTypeString).
			WithUsage("Server").
			WithGroup("Connect").
			WithDefault("tcp://localhost:1883").
			WithEditable(true),
		config.NewVariable(mqtt.ConfigConnectionTimeout, config.ValueTypeDuration).
			WithUsage("Timeout").
			WithGroup("Connect").
			WithDefault(opts.ConnectTimeout),
		config.NewVariable(mqtt.ConfigConnectionRetryEnabled, config.ValueTypeBool).
			WithUsage("Retry").
			WithGroup("Connect").
			WithDefault(opts.ConnectRetry),
		config.NewVariable(mqtt.ConfigConnectionRetryInterval, config.ValueTypeDuration).
			WithUsage("Sets the time that will be waited between connection attempts when initially connecting if " + mqtt.ConfigConnectionRetryEnabled + " is true").
			WithGroup("Connect").
			WithDefault(opts.ConnectRetryInterval),
		config.NewVariable(mqtt.ConfigKeepAlive, config.ValueTypeDuration).
			WithUsage("Will set the amount of time (in seconds) that the client should wait before sending a PING request to the broker. This will allow the client to know that a connection has not been lost with the server.").
			WithGroup("Connect").
			WithDefault(time.Duration(opts.KeepAlive) * time.Second),
		config.NewVariable(mqtt.ConfigWriteTimeout, config.ValueTypeDuration).
			WithUsage("Limit on how long a mqtt publish should block until it unblocks with a timeout error. A duration of 0 never times out.").
			WithGroup("Connect").
			WithDefault(opts.WriteTimeout),
		config.NewVariable(mqtt.ConfigPingTimeout, config.ValueTypeDuration).
			WithUsage("Will set the amount of time (in seconds) that the client will wait after sending a PING request to the broker, before deciding that the connection has been lost.").
			WithGroup("Connect").
			WithDefault(opts.PingTimeout),
		config.NewVariable(mqtt.ConfigTokenWaitTimeout, config.ValueTypeDuration).
			WithUsage("Token wait timeout for requests.").
			WithGroup("Connect").
			WithDefault(time.Second * 5).
			WithEditable(true),
		config.NewVariable(mqtt.ConfigReconnectEnabled, config.ValueTypeBool).
			WithUsage("Auto reconnect enabled.").
			WithGroup("Reconnect").
			WithDefault(opts.AutoReconnect),
		config.NewVariable(mqtt.ConfigReconnectMaxInterval, config.ValueTypeDuration).
			WithUsage("Reconnect max interval.").
			WithGroup("Reconnect").
			WithDefault(opts.MaxReconnectInterval),
		config.NewVariable(mqtt.ConfigLWTEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Last Will and Testament").
			WithDefault(opts.WillEnabled),
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
			WithDefault(opts.WillQos).
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
			WithDefault(opts.WillRetained),
		config.NewVariable(mqtt.ConfigPayloadCacheEnabled, config.ValueTypeBool).
			WithUsage("Payload cache enabled").
			WithGroup("Cache").
			WithEditable(true).
			WithDefault(true),
		config.NewVariable(mqtt.ConfigPayloadCacheSize, config.ValueTypeInt64).
			WithUsage("Payload cache size").
			WithGroup("Cache").
			WithEditable(true).
			WithDefault(1000),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{
			mqtt.ConfigServers,
			mqtt.ConfigClientID,
			mqtt.ConfigUsername,
			mqtt.ConfigPassword,
		}, c.watchConnect),
		config.NewWatcher([]string{
			mqtt.ConfigPayloadCacheEnabled,
		}, c.watchPayloadCacheEnabled),
		config.NewWatcher([]string{
			mqtt.ConfigPayloadCacheSize,
		}, c.watchPayloadCacheSize),
	}
}

func (c *Component) watchConnect(_ string, _ interface{}, _ interface{}) {
	if err := c.initClient(); err == nil {
		if err := c.initSubscribers(); err != nil {
			c.logger.Warn("Failed init MQTT subscribers", "error", err.Error())
		}
	} else {
		c.logger.Warn("Failed init MQTT client", "error", err.Error())
	}
}

func (c *Component) watchPayloadCacheEnabled(_ string, newValue interface{}, _ interface{}) {
	if !newValue.(bool) {
		c.payloadCache.Purge()
	}
}

func (c *Component) watchPayloadCacheSize(_ string, newValue interface{}, _ interface{}) {
	err := c.payloadCache.Resize(newValue.(int))

	if err != nil {
		c.logger.Error("Failed resize payload cache", "error", err.Error())
	}
}
