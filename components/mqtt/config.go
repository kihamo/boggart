package mqtt

const (
	DefaultConnectionAttempts = 3

	ConfigServers            = ComponentName + ".servers"
	ConfigClientID           = ComponentName + ".client-id"
	ConfigUsername           = ComponentName + ".username"
	ConfigPassword           = ComponentName + ".password"
	ConfigConnectionAttempts = ComponentName + ".connection.attempts"
	ConfigConnectionTimeout  = ComponentName + ".connection.timeout"
	ConfigClearSession       = ComponentName + ".clear-session"
	ConfigResumeSubs         = ComponentName + ".resume-subs"
	ConfigWriteTimeout       = ComponentName + ".write-timeout"
	ConfigLWTEnabled         = ComponentName + ".lwt.enabled"
	ConfigLWTTopic           = ComponentName + ".lwt.topic"
	ConfigLWTPayload         = ComponentName + ".lwt.payload"
	ConfigLWTQOS             = ComponentName + ".lwt.qos"
	ConfigLWTRetained        = ComponentName + ".lwt.retained"
)
