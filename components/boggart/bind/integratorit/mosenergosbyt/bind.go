package mosenergosbyt

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/integratorit/mosenergosbyt"
)

var services = map[string]string{
	"взнос на капитальный ремонт": "10001",
	"обращение с тко":             "10002",
}

type Bind struct {
	boggart.BindBase
	di.MQTTBind
	di.WorkersBind

	config *Config
	client *mosenergosbyt.Client
}
