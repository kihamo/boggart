package voice

const (
	ComponentName    = "voice"
	ComponentVersion = "0.1.0"

	MQTTTopicSimpleText   = ComponentName + "/speech/text"
	MQTTTopicJSONText     = ComponentName + "/speech/json"
	MQTTTopicPlayerURL    = ComponentName + "/player/url"
	MQTTTopicPlayerPause  = ComponentName + "/player/pause"
	MQTTTopicPlayerStop   = ComponentName + "/player/stop"
	MQTTTopicPlayerPlay   = ComponentName + "/player/play"
	MQTTTopicPlayerStatus = ComponentName + "/player/status"
)
