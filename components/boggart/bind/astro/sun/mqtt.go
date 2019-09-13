package sun

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	return []mqtt.Topic{
		b.config.TopicNadir,
		b.config.TopicNightBeforeStart,
		b.config.TopicNightBeforeEnd,
		b.config.TopicNightBeforeDuration,
		b.config.TopicAstronomicalDawnStart,
		b.config.TopicAstronomicalDawnEnd,
		b.config.TopicAstronomicalDawnDuration,
		b.config.TopicNauticalDawnStart,
		b.config.TopicNauticalDawnEnd,
		b.config.TopicNauticalDawnDuration,
		b.config.TopicCivilDawnStart,
		b.config.TopicCivilDawnEnd,
		b.config.TopicCivilDawnDuration,
		b.config.TopicRiseStart,
		b.config.TopicRiseEnd,
		b.config.TopicRiseDuration,
		b.config.TopicSolarNoon,
		b.config.TopicSetStart,
		b.config.TopicSetEnd,
		b.config.TopicSetDuration,
		b.config.TopicCivilDuskStart,
		b.config.TopicCivilDuskEnd,
		b.config.TopicCivilDuskDuration,
		b.config.TopicNauticalDuskStart,
		b.config.TopicNauticalDuskEnd,
		b.config.TopicNauticalDuskDuration,
		b.config.TopicAstronomicalDuskStart,
		b.config.TopicAstronomicalDuskEnd,
		b.config.TopicAstronomicalDuskDuration,
		b.config.TopicNightAfterStart,
		b.config.TopicNightAfterEnd,
		b.config.TopicNightAfterDuration,
	}
}
