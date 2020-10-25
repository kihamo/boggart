package pass24online

import (
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/pass24online"
)

type Bind struct {
	di.LoggerBind
	di.MQTTBind
	di.WidgetBind
	di.WorkersBind

	config   *Config
	provider *pass24online.Client

	feedStartDatetime *atomic.Time
}

func (b *Bind) Run() error {
	b.feedStartDatetime.Set(time.Now().Add(time.Hour * -24))

	return nil
}
