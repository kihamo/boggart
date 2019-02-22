package boggart

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config
}

func (b *Bind) Run() error {
	b.UpdateStatus(boggart.BindStatusOnline)

	ctx := context.Background()
	sn := mqtt.NameReplace(b.config.ApplicationName)

	b.MQTTPublishAsync(ctx, MQTTPublishTopicApplicationName.Format(sn), b.config.ApplicationName)
	b.MQTTPublishAsync(ctx, MQTTPublishTopicApplicationVersion.Format(sn), b.config.ApplicationVersion)
	b.MQTTPublishAsync(ctx, MQTTPublishTopicApplicationBuild.Format(sn), b.config.ApplicationBuild)

	return nil
}
