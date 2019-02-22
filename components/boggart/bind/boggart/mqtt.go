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
	sn := mqtt.NameReplace(b.config.ApplicationName)

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicApplicationName.Format(sn)),
		mqtt.Topic(MQTTPublishTopicApplicationVersion.Format(sn)),
		mqtt.Topic(MQTTPublishTopicApplicationBuild.Format(sn)),
	}
}
