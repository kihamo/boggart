package sun

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)
	config.TopicNadir = config.TopicNadir.Format(config.Name)
	config.TopicNightBeforeStart = config.TopicNightBeforeStart.Format(config.Name)
	config.TopicNightBeforeEnd = config.TopicNightBeforeEnd.Format(config.Name)
	config.TopicNightBeforeDuration = config.TopicNightBeforeDuration.Format(config.Name)
	config.TopicAstronomicalDawnStart = config.TopicAstronomicalDawnStart.Format(config.Name)
	config.TopicAstronomicalDawnEnd = config.TopicAstronomicalDawnEnd.Format(config.Name)
	config.TopicAstronomicalDawnDuration = config.TopicAstronomicalDawnDuration.Format(config.Name)
	config.TopicNauticalDawnStart = config.TopicNauticalDawnStart.Format(config.Name)
	config.TopicNauticalDawnEnd = config.TopicNauticalDawnEnd.Format(config.Name)
	config.TopicNauticalDawnDuration = config.TopicNauticalDawnDuration.Format(config.Name)
	config.TopicCivilDawnStart = config.TopicCivilDawnStart.Format(config.Name)
	config.TopicCivilDawnEnd = config.TopicCivilDawnEnd.Format(config.Name)
	config.TopicCivilDawnDuration = config.TopicCivilDawnDuration.Format(config.Name)
	config.TopicRiseStart = config.TopicRiseStart.Format(config.Name)
	config.TopicRiseEnd = config.TopicRiseEnd.Format(config.Name)
	config.TopicRiseDuration = config.TopicRiseDuration.Format(config.Name)
	config.TopicSolarNoon = config.TopicSolarNoon.Format(config.Name)
	config.TopicSetStart = config.TopicSetStart.Format(config.Name)
	config.TopicSetEnd = config.TopicSetEnd.Format(config.Name)
	config.TopicSetDuration = config.TopicSetDuration.Format(config.Name)
	config.TopicCivilDuskStart = config.TopicCivilDuskStart.Format(config.Name)
	config.TopicCivilDuskEnd = config.TopicCivilDuskEnd.Format(config.Name)
	config.TopicCivilDuskDuration = config.TopicCivilDuskDuration.Format(config.Name)
	config.TopicNauticalDuskStart = config.TopicNauticalDuskStart.Format(config.Name)
	config.TopicNauticalDuskEnd = config.TopicNauticalDuskEnd.Format(config.Name)
	config.TopicNauticalDuskDuration = config.TopicNauticalDuskDuration.Format(config.Name)
	config.TopicAstronomicalDuskStart = config.TopicAstronomicalDuskStart.Format(config.Name)
	config.TopicAstronomicalDuskEnd = config.TopicAstronomicalDuskEnd.Format(config.Name)
	config.TopicAstronomicalDuskDuration = config.TopicAstronomicalDuskDuration.Format(config.Name)
	config.TopicNightAfterStart = config.TopicNightAfterStart.Format(config.Name)
	config.TopicNightAfterEnd = config.TopicNightAfterEnd.Format(config.Name)
	config.TopicNightAfterDuration = config.TopicNightAfterDuration.Format(config.Name)

	bind := &Bind{
		config: config,
	}

	return bind, nil
}
