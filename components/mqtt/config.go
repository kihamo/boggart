package mqtt

const (
	ConfigServers             = ComponentName + ".servers"
	ConfigClientID            = ComponentName + ".client-id"
	ConfigUsername            = ComponentName + ".username"
	ConfigPassword            = ComponentName + ".password"
	ConfigConnectionTimeout   = ComponentName + ".connection.timeout"
	ConfigClearSession        = ComponentName + ".clear-session"
	ConfigResumeSubs          = ComponentName + ".resume-subs"
	ConfigWriteTimeout        = ComponentName + ".write-timeout"
	ConfigLWTEnabled          = ComponentName + ".lwt.enabled"
	ConfigLWTTopic            = ComponentName + ".lwt.topic"
	ConfigLWTPayload          = ComponentName + ".lwt.payload"
	ConfigLWTQOS              = ComponentName + ".lwt.qos"
	ConfigLWTRetained         = ComponentName + ".lwt.retained"
	ConfigPayloadCacheEnabled = ComponentName + ".cache.payload.enabled"
	ConfigPayloadCacheSize    = ComponentName + ".cache.payload.size"
)
