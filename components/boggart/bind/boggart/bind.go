package boggart

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/shadow"
	"go.uber.org/multierr"
)

type Bind struct {
	di.ConfigBind
	di.MetaBind
	di.MQTTBind
	di.WidgetBind

	application shadow.Application
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() (err error) {
	ctx := context.Background()
	cfg := b.config()

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicName.Format(cfg.ApplicationName), cfg.ApplicationName); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicVersion.Format(cfg.ApplicationName), cfg.ApplicationVersion); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicBuild.Format(cfg.ApplicationName), cfg.ApplicationBuild); e != nil {
		err = multierr.Append(err, e)
	}

	// защита на случай, если запустят с retained
	if e := b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicShutdown.Format(cfg.ApplicationName), false); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}

func (b *Bind) SetApplication(a shadow.Application) {
	b.application = a
}
