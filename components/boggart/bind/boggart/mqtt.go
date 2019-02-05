package boggart

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/boggart/+"

	MQTTPublishTopicApplicationName    = MQTTPrefix + "/application/name"
	MQTTPublishTopicApplicationVersion = MQTTPrefix + "/application/version"
	MQTTPublishTopicApplicationBuild   = MQTTPrefix + "/application/build"
)

func (b *Bind) SetMQTTClient(client mqtt.Component) {
	b.BindMQTT.SetMQTTClient(client)

	if client != nil {
		ctx := context.Background()
		sn := mqtt.NameReplace(b.config.ApplicationName)

		client.PublishAsync(ctx, MQTTPublishTopicApplicationName.Format(sn), 2, true, b.config.ApplicationName)
		client.PublishAsync(ctx, MQTTPublishTopicApplicationVersion.Format(sn), 2, true, b.config.ApplicationVersion)
		client.PublishAsync(ctx, MQTTPublishTopicApplicationBuild.Format(sn), 2, true, b.config.ApplicationBuild)
	}
}

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(b.config.ApplicationName)

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicApplicationName.Format(sn)),
		mqtt.Topic(MQTTPublishTopicApplicationVersion.Format(sn)),
		mqtt.Topic(MQTTPublishTopicApplicationBuild.Format(sn)),
	}
}
