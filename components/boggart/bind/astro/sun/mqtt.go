package sun

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/astro/sun/+/"

	MQTTPublishTopicNadir                    = MQTTPrefix + "nadir"
	MQTTPublishTopicNightBeforeStart         = MQTTPrefix + "night/before/start"
	MQTTPublishTopicNightBeforeEnd           = MQTTPrefix + "night/before/end"
	MQTTPublishTopicNightBeforeDuration      = MQTTPrefix + "night/before/duration"
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
	MQTTPublishTopicNightAfterStart          = MQTTPrefix + "night/after/start"
	MQTTPublishTopicNightAfterEnd            = MQTTPrefix + "night/after/end"
	MQTTPublishTopicNightAfterDuration       = MQTTPrefix + "night/after/duration"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := b.SerialNumber()

	return []mqtt.Topic{
		MQTTPublishTopicNadir.Format(sn),
		MQTTPublishTopicNightBeforeStart.Format(sn),
		MQTTPublishTopicNightBeforeEnd.Format(sn),
		MQTTPublishTopicNightBeforeDuration.Format(sn),
		MQTTPublishTopicAstronomicalDawnStart.Format(sn),
		MQTTPublishTopicAstronomicalDawnEnd.Format(sn),
		MQTTPublishTopicAstronomicalDawnDuration.Format(sn),
		MQTTPublishTopicNauticalDawnStart.Format(sn),
		MQTTPublishTopicNauticalDawnEnd.Format(sn),
		MQTTPublishTopicNauticalDawnDuration.Format(sn),
		MQTTPublishTopicCivilDawnStart.Format(sn),
		MQTTPublishTopicCivilDawnEnd.Format(sn),
		MQTTPublishTopicCivilDawnDuration.Format(sn),
		MQTTPublishTopicRiseStart.Format(sn),
		MQTTPublishTopicRiseEnd.Format(sn),
		MQTTPublishTopicRiseDuration.Format(sn),
		MQTTPublishTopicSolarNoon.Format(sn),
		MQTTPublishTopicSetStart.Format(sn),
		MQTTPublishTopicSetEnd.Format(sn),
		MQTTPublishTopicSetDuration.Format(sn),
		MQTTPublishTopicCivilDuskStart.Format(sn),
		MQTTPublishTopicCivilDuskEnd.Format(sn),
		MQTTPublishTopicCivilDuskDuration.Format(sn),
		MQTTPublishTopicNauticalDuskStart.Format(sn),
		MQTTPublishTopicNauticalDuskEnd.Format(sn),
		MQTTPublishTopicNauticalDuskDuration.Format(sn),
		MQTTPublishTopicAstronomicalDuskStart.Format(sn),
		MQTTPublishTopicAstronomicalDuskEnd.Format(sn),
		MQTTPublishTopicAstronomicalDuskDuration.Format(sn),
		MQTTPublishTopicNightAfterStart.Format(sn),
		MQTTPublishTopicNightAfterEnd.Format(sn),
		MQTTPublishTopicNightAfterDuration.Format(sn),
	}
}
