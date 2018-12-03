package mqtt

const (
	DefaultConnectionAttempts = 3

	ConfigServers            = ComponentName + ".servers"
	ConfigUsername           = ComponentName + ".username"
	ConfigPassword           = ComponentName + ".password"
	ConfigConnectionAttempts = ComponentName + ".connection.attempts"
	ConfigConnectionTimeout  = ComponentName + ".connection.timeout"
)
