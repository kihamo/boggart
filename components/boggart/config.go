package boggart

const (
	ConfigListenerTelegramChats       = ComponentName + ".listener.telegram.chats"
	ConfigRS485Address                = ComponentName + ".rs485.path"
	ConfigRS485Timeout                = ComponentName + ".rs485.timeout"
	ConfigGPIOPins                    = ComponentName + ".gpio.pins"
	ConfigMercuryRepeatInterval       = ComponentName + ".mercury.repeat-interval"
	ConfigMercuryDeviceAddress        = ComponentName + ".mercury.device-address"
	ConfigPulsarEnabled               = ComponentName + ".pulsar.enabled"
	ConfigPulsarRepeatInterval        = ComponentName + ".pulsar.repeat-interval"
	ConfigPulsarHeatMeterAddress      = ComponentName + ".pulsar.heat-meter.address"
	ConfigPulsarColdWaterSerialNumber = ComponentName + ".pulsar.cold-water.serial-number"
	ConfigPulsarColdWaterPulseInput   = ComponentName + ".pulsar.cold-water.pulse-input"
	ConfigPulsarColdWaterStartValue   = ComponentName + ".pulsar.cold-water.start-value"
	ConfigPulsarHotWaterSerialNumber  = ComponentName + ".pulsar.hot-water.serial-number"
	ConfigPulsarHotWaterPulseInput    = ComponentName + ".pulsar.hot-water.pulse-input"
	ConfigPulsarHotWaterStartValue    = ComponentName + ".pulsar.hot-water.start-value"
	ConfigMQTTOwnTracksEnabled        = ComponentName + ".mqtt.own-tracks.enabled"
	ConfigMQTTWOLEnabled              = ComponentName + ".mqtt.wol.enabled"
	ConfigMQTTAnnotationsEnabled      = ComponentName + ".mqtt.annotations.enabled"
	ConfigMQTTMessengersEnabled       = ComponentName + ".mqtt.messengers.enabled"
	ConfigConfigYAML                  = ComponentName + ".config.yaml"
)
