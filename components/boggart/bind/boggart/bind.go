package boggart

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"go.uber.org/multierr"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config
}

func (b *Bind) Run() (err error) {
	b.UpdateStatus(boggart.BindStatusOnline)

	ctx := context.Background()

	if e := b.MQTTPublishAsync(ctx, b.config.TopicName, b.config.ApplicationName); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, b.config.TopicVersion, b.config.ApplicationVersion); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, b.config.TopicBuild, b.config.ApplicationBuild); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}
