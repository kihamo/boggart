package pass24online

import (
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/pass24online"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.WidgetBind
	di.WorkersBind

	provider *pass24online.Client

	feedStartDatetime *atomic.Time
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	b.feedStartDatetime.Set(time.Now().Add(time.Hour * -24))
	b.provider = pass24online.New(cfg.Phone, cfg.Password, cfg.Debug, swagger.NewLogger(
		func(message string) {
			b.Logger().Info(message)
		},
		func(message string) {
			b.Logger().Debug(message)
		}))

	return nil
}
