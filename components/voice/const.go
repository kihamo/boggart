package voice

const (
	ComponentName    = "voice"
	ComponentVersion = "0.1.0"

	MQTTTopicSimpleText      = ComponentName + "/speech/text"
	MQTTTopicJSONText        = ComponentName + "/speech/json"
	MQTTTopicPlayerURL       = ComponentName + "/player/url"
	MQTTTopicPlayerStatus    = ComponentName + "/player/status"
	MQTTTopicPlayerAction    = ComponentName + "/player/action"
	MQTTTopicPlayerVolume    = ComponentName + "/player/volume"
	MQTTTopicPlayerVolumeSet = ComponentName + "/player/volume/value"
)
