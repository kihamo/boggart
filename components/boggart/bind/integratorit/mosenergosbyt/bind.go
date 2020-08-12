package mosenergosbyt

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/integratorit/mosenergosbyt"
)

var services = map[string]string{
	"взнос на капитальный ремонт": "10001",
	"обращение с тко":             "10002",
}

const (
	layoutPeriod = "2006-01-02"
)

type Bind struct {
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind
	di.WidgetBind

	config *Config
	client *mosenergosbyt.Client
}
