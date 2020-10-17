package boggart

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/shadow"
	"go.uber.org/multierr"
)

type Bind struct {
	di.MQTTBind

	config      *Config
	application shadow.Application
}

func (b *Bind) Run() (err error) {
	ctx := context.Background()

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicName, b.config.ApplicationName); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicVersion, b.config.ApplicationVersion); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicBuild, b.config.ApplicationBuild); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}

func (b *Bind) SetApplication(a shadow.Application) {
	b.application = a
}
