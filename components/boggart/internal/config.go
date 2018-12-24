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
		config.NewVariable(boggart.ConfigDevicesManagerCheckInterval, config.ValueTypeDuration).
			WithUsage("Health check interval").
			WithGroup("Device manager").
			WithEditable(true).
			WithDefault(time.Minute),
		config.NewVariable(boggart.ConfigDevicesManagerCheckTimeout, config.ValueTypeDuration).
			WithUsage("Health check timeout").
			WithGroup("Device manager").
			WithEditable(true).
			WithDefault(DefaultTimeoutChecker),
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
		config.NewVariable(boggart.ConfigGPIOPins, config.ValueTypeString).
			WithUsage("Pins listener").
			WithGroup("GPIO"),
		config.NewVariable(boggart.ConfigCameraHikVisionAddresses, config.ValueTypeString).
			WithUsage("Addresses").
			WithGroup("Cameras").
			WithView([]string{config.ViewTags}),
		config.NewVariable(boggart.ConfigCameraHikVisionRepeatInterval, config.ValueTypeDuration).
			WithUsage("Repeat interval").
			WithGroup("Cameras").
			WithDefault(time.Minute),
		config.NewVariable(boggart.ConfigMercuryRepeatInterval, config.ValueTypeDuration).
			WithUsage("Repeat interval").
			WithGroup("Mercury devices").
			WithDefault(time.Minute * 2),
		config.NewVariable(boggart.ConfigMercuryDeviceAddress, config.ValueTypeString).
			WithUsage("Device address in format XXXXXX (last 6 digits of device serial number)").
			WithGroup("Mercury devices"),
		config.NewVariable(boggart.ConfigMikrotikRepeatInterval, config.ValueTypeDuration).
			WithUsage("Repeat interval").
			WithGroup("Mikrotik devices").
			WithDefault(time.Minute * 5),
		config.NewVariable(boggart.ConfigMikrotikAddresses, config.ValueTypeString).
			WithUsage("API address in format host:port").
			WithGroup("Mikrotik devices").
			WithView([]string{config.ViewTags}),
		config.NewVariable(boggart.ConfigMikrotikTimeout, config.ValueTypeDuration).
			WithUsage("Request timeout").
			WithGroup("Mikrotik devices").
			WithDefault(time.Second * 10),
		config.NewVariable(boggart.ConfigMobileEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Mobile accounts"),
		config.NewVariable(boggart.ConfigMobileRepeatInterval, config.ValueTypeDuration).
			WithUsage("Repeat interval").
			WithGroup("Mobile accounts").
			WithDefault(time.Minute * 30),
		config.NewVariable(boggart.ConfigMobileMegafonPhone, config.ValueTypeString).
			WithUsage("Phone number").
			WithGroup("Mobile Megafon accounts"),
		config.NewVariable(boggart.ConfigMobileMegafonPassword, config.ValueTypeString).
			WithUsage("Password").
			WithGroup("Mobile Megafon accounts").
			WithView([]string{config.ViewPassword}),
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
		config.NewVariable(boggart.ConfigSoftVideoEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("SoftVideo provider"),
		config.NewVariable(boggart.ConfigSoftVideoRepeatInterval, config.ValueTypeDuration).
			WithUsage("Repeat interval").
			WithGroup("SoftVideo provider").
			WithDefault(time.Hour * 12),
		config.NewVariable(boggart.ConfigSoftVideoLogin, config.ValueTypeString).
			WithUsage("Login").
			WithGroup("SoftVideo provider"),
		config.NewVariable(boggart.ConfigSoftVideoPassword, config.ValueTypeString).
			WithUsage("Password").
			WithGroup("SoftVideo provider").
			WithView([]string{config.ViewPassword}),
		config.NewVariable(boggart.ConfigSensorBME280Enabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Sensor 280"),
		config.NewVariable(boggart.ConfigSensorBME280RepeatInterval, config.ValueTypeDuration).
			WithUsage("Repeat interval").
			WithGroup("Sensor 280").
			WithDefault(time.Minute),
		config.NewVariable(boggart.ConfigSensorBME280Bus, config.ValueTypeInt).
			WithUsage("Bus number").
			WithGroup("Sensor 280").
			WithDefault(1),
		config.NewVariable(boggart.ConfigSensorBME280Address, config.ValueTypeInt).
			WithUsage("Address").
			WithGroup("Sensor 280").
			WithDefault(0x76),
		config.NewVariable(boggart.ConfigOwnTracksEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("OwnTracks"),
		config.NewVariable(boggart.ConfigWOLEnabled, config.ValueTypeBool).
			WithUsage("Enabled").
			WithGroup("Wake-on-LAN"),
		config.NewVariable(boggart.ConfigSocketsBroadlink, config.ValueTypeString).
			WithUsage("Address of Broadlink in format ip:mac").
			WithGroup("Sockets").
			WithView([]string{config.ViewTags}),
		config.NewVariable(boggart.ConfigRemoteControlBroadlink, config.ValueTypeString).
			WithUsage("Address of Broadlink in format ip:mac").
			WithGroup("Remote control").
			WithView([]string{config.ViewTags}),
		config.NewVariable(boggart.ConfigLEDWiFi, config.ValueTypeString).
			WithUsage("Address of WiFi LED in format hostname or ip without port").
			WithGroup("LED control").
			WithView([]string{config.ViewTags}),
		config.NewVariable(boggart.ConfigTVLGWebOS, config.ValueTypeString).
			WithUsage("Address of LG on WebOS in format hostname or ip without port and register key").
			WithGroup("TV").
			WithView([]string{config.ViewTags}),
		config.NewVariable(boggart.ConfigTVSamsung, config.ValueTypeString).
			WithUsage("Address of Samsung in format hostname").
			WithGroup("TV").
			WithView([]string{config.ViewTags}),
	}
}

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{
			boggart.ConfigDevicesManagerCheckInterval,
		}, c.watchDevicesManagerCheckInterval),
		config.NewWatcher([]string{
			boggart.ConfigDevicesManagerCheckTimeout,
		}, c.watchDevicesManagerCheckTimeout),
	}
}

func (c *Component) watchDevicesManagerCheckInterval(_ string, newValue interface{}, _ interface{}) {
	c.devicesManager.SetCheckerTickerDuration(newValue.(time.Duration))
}

func (c *Component) watchDevicesManagerCheckTimeout(_ string, newValue interface{}, _ interface{}) {
	c.devicesManager.SetCheckerTimeout(newValue.(time.Duration))
}
