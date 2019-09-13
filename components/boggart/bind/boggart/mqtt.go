package boggart

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/boggart/+"

	MQTTPublishTopicApplicationName    = MQTTPrefix + "/application/name"
	MQTTPublishTopicApplicationVersion = MQTTPrefix + "/application/version"
	MQTTPublishTopicApplicationBuild   = MQTTPrefix + "/application/build"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		MQTTPublishTopicApplicationName.Format(b.config.ApplicationName),
		MQTTPublishTopicApplicationVersion.Format(b.config.ApplicationName),
		MQTTPublishTopicApplicationBuild.Format(b.config.ApplicationName),
	}
}
