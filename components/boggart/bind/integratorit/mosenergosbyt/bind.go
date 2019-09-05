package mosenergosbyt

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/integratorit/mosenergosbyt"
)

var services = map[string]string{
	"взнос на капитальный ремонт": "10001",
	"обращение с тко":             "10002",
}

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	client *mosenergosbyt.Client
	config *Config
}
