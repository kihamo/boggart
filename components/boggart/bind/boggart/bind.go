package boggart

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
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
	sn := mqtt.NameReplace(b.config.ApplicationName)

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicApplicationName.Format(sn), b.config.ApplicationName); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicApplicationVersion.Format(sn), b.config.ApplicationVersion); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTTPublishAsync(ctx, MQTTPublishTopicApplicationBuild.Format(sn), b.config.ApplicationBuild); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}
