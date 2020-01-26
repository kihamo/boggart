package samsung_tizen

import (
	"sync"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/samsung/tv"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.WorkersBind
	di.LoggerBind
	di.ProbesBind

	config *Config
	mutex  sync.RWMutex
	client *tv.ApiV2
	mac    string
}
