package sun

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/astro/sun/+/"

	MQTTPublishTopicNadir                    = MQTTPrefix + "nadir"
	MQTTPublishTopicNightEnd                 = MQTTPrefix + "night/end"
	MQTTPublishTopicAstronomicalDawnStart    = MQTTPrefix + "dawn/astronomical/start"
	MQTTPublishTopicAstronomicalDawnEnd      = MQTTPrefix + "dawn/astronomical/end"
	MQTTPublishTopicAstronomicalDawnDuration = MQTTPrefix + "dawn/astronomical/duration"
	MQTTPublishTopicNauticalDawnStart        = MQTTPrefix + "dawn/nautical/start"
	MQTTPublishTopicNauticalDawnEnd          = MQTTPrefix + "dawn/nautical/end"
	MQTTPublishTopicNauticalDawnDuration     = MQTTPrefix + "dawn/nautical/duration"
	MQTTPublishTopicCivilDawnStart           = MQTTPrefix + "dawn/civil/start"
	MQTTPublishTopicCivilDawnEnd             = MQTTPrefix + "dawn/civil/end"
	MQTTPublishTopicCivilDawnDuration        = MQTTPrefix + "dawn/civil/duration"
	MQTTPublishTopicRiseStart                = MQTTPrefix + "rise/start"
	MQTTPublishTopicRiseEnd                  = MQTTPrefix + "rise/end"
	MQTTPublishTopicRiseDuration             = MQTTPrefix + "rise/duration"
	MQTTPublishTopicSolarNoon                = MQTTPrefix + "solar-noon"
	MQTTPublishTopicSetStart                 = MQTTPrefix + "set/start"
	MQTTPublishTopicSetEnd                   = MQTTPrefix + "set/end"
	MQTTPublishTopicSetDuration              = MQTTPrefix + "set/duration"
	MQTTPublishTopicCivilDuskStart           = MQTTPrefix + "dusk/civil/start"
	MQTTPublishTopicCivilDuskEnd             = MQTTPrefix + "dusk/civil/end"
	MQTTPublishTopicCivilDuskDuration        = MQTTPrefix + "dusk/civil/duration"
	MQTTPublishTopicNauticalDuskStart        = MQTTPrefix + "dusk/nautical/start"
	MQTTPublishTopicNauticalDuskEnd          = MQTTPrefix + "dusk/nautical/end"
	MQTTPublishTopicNauticalDuskDuration     = MQTTPrefix + "dusk/nautical/duration"
	MQTTPublishTopicAstronomicalDuskStart    = MQTTPrefix + "dusk/astronomical/start"
	MQTTPublishTopicAstronomicalDuskEnd      = MQTTPrefix + "dusk/astronomical/end"
	MQTTPublishTopicAstronomicalDuskDuration = MQTTPrefix + "dusk/astronomical/duration"
	MQTTPublishTopicNightStart               = MQTTPrefix + "night/start"
	MQTTPublishTopicNightDuration            = MQTTPrefix + "night/duration"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicNadir.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNightEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicAstronomicalDawnStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicAstronomicalDawnEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicAstronomicalDawnDuration.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNauticalDawnStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNauticalDawnEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNauticalDawnDuration.Format(sn)),
		mqtt.Topic(MQTTPublishTopicCivilDawnStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicCivilDawnEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicCivilDawnDuration.Format(sn)),
		mqtt.Topic(MQTTPublishTopicRiseStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicRiseEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicRiseDuration.Format(sn)),
		mqtt.Topic(MQTTPublishTopicSolarNoon.Format(sn)),
		mqtt.Topic(MQTTPublishTopicSetStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicSetEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicSetDuration.Format(sn)),
		mqtt.Topic(MQTTPublishTopicCivilDuskStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicCivilDuskEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicCivilDuskDuration.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNauticalDuskStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNauticalDuskEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNauticalDuskDuration.Format(sn)),
		mqtt.Topic(MQTTPublishTopicAstronomicalDuskStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicAstronomicalDuskEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicAstronomicalDuskDuration.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNightStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNightDuration.Format(sn)),
	}
}
