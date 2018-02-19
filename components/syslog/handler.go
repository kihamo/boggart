package syslog

type HasHandler interface {
	SyslogHandler(map[string]interface{})
}
