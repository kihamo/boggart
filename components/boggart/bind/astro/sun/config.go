package sun

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	Lat float64 `valid:"required"`
	Lon float64 `valid:"required"`

	TopicNadir                    mqtt.Topic `mapstructure:"topic_nadir" yaml:"topic_nadir"`
	TopicNightBeforeStart         mqtt.Topic `mapstructure:"topic_night_before_start" yaml:"topic_night_before_start"`
	TopicNightBeforeEnd           mqtt.Topic `mapstructure:"topic_night_before_end" yaml:"topic_night_before_end"`
	TopicNightBeforeDuration      mqtt.Topic `mapstructure:"topic_night_before_duration" yaml:"topic_night_before_duration"`
	TopicAstronomicalDawnStart    mqtt.Topic `mapstructure:"topic_astronomical_dawn_start" yaml:"topic_astronomical_dawn_start"`
	TopicAstronomicalDawnEnd      mqtt.Topic `mapstructure:"topic_astronomical_dawn_end" yaml:"topic_astronomical_dawn_end"`
	TopicAstronomicalDawnDuration mqtt.Topic `mapstructure:"topic_astronomical_dawn_duration" yaml:"topic_astronomical_dawn_duration"`
	TopicNauticalDawnStart        mqtt.Topic `mapstructure:"topic_nautical_dawn_start" yaml:"topic_nautical_dawn_start"`
	TopicNauticalDawnEnd          mqtt.Topic `mapstructure:"topic_nautical_dawn_end" yaml:"topic_nautical_dawn_end"`
	TopicNauticalDawnDuration     mqtt.Topic `mapstructure:"topic_nautical_dawn_duration" yaml:"topic_nautical_dawn_duration"`
	TopicCivilDawnStart           mqtt.Topic `mapstructure:"topic_civil_dawn_start" yaml:"topic_civil_dawn_start"`
	TopicCivilDawnEnd             mqtt.Topic `mapstructure:"topic_civil_dawn_end" yaml:"topic_civil_dawn_end"`
	TopicCivilDawnDuration        mqtt.Topic `mapstructure:"topic_civil_dawn_duration" yaml:"topic_civil_dawn_duration"`
	TopicRiseStart                mqtt.Topic `mapstructure:"topic_rise_start" yaml:"topic_rise_start"`
	TopicRiseEnd                  mqtt.Topic `mapstructure:"topic_rise_end" yaml:"topic_rise_end"`
	TopicRiseDuration             mqtt.Topic `mapstructure:"topic_rise_duration" yaml:"topic_rise_duration"`
	TopicSolarNoon                mqtt.Topic `mapstructure:"topic_solar_noon" yaml:"topic_solar_noon"`
	TopicSetStart                 mqtt.Topic `mapstructure:"topic_set_start" yaml:"topic_set_start"`
	TopicSetEnd                   mqtt.Topic `mapstructure:"topic_set_end" yaml:"topic_set_end"`
	TopicSetDuration              mqtt.Topic `mapstructure:"topic_set_duration" yaml:"topic_set_duration"`
	TopicCivilDuskStart           mqtt.Topic `mapstructure:"topic_civil_dusk_start" yaml:"topic_civil_dusk_start"`
	TopicCivilDuskEnd             mqtt.Topic `mapstructure:"topic_civil_dusk_end" yaml:"topic_civil_dusk_end"`
	TopicCivilDuskDuration        mqtt.Topic `mapstructure:"topic_civil_dusk_duration" yaml:"topic_civil_dusk_duration"`
	TopicNauticalDuskStart        mqtt.Topic `mapstructure:"topic_nautical_dusk_start" yaml:"topic_nautical_dusk_start"`
	TopicNauticalDuskEnd          mqtt.Topic `mapstructure:"topic_nautical_dusk_end" yaml:"topic_nautical_dusk_end"`
	TopicNauticalDuskDuration     mqtt.Topic `mapstructure:"topic_nautical_dusk_duration" yaml:"topic_nautical_dusk_duration"`
	TopicAstronomicalDuskStart    mqtt.Topic `mapstructure:"topic_astronomical_dusk_start" yaml:"topic_astronomical_dusk_start"`
	TopicAstronomicalDuskEnd      mqtt.Topic `mapstructure:"topic_astronomical_dusk_end" yaml:"topic_astronomical_dusk_end"`
	TopicAstronomicalDuskDuration mqtt.Topic `mapstructure:"topic_astronomical_dusk_duration" yaml:"topic_astronomical_dusk_duration"`
	TopicNightAfterStart          mqtt.Topic `mapstructure:"topic_night_after_start" yaml:"topic_night_after_start"`
	TopicNightAfterEnd            mqtt.Topic `mapstructure:"topic_night_after_end" yaml:"topic_night_after_end"`
	TopicNightAfterDuration       mqtt.Topic `mapstructure:"topic_night_after_duration" yaml:"topic_night_after_duration"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/astro/sun/+/"

	return &Config{
		TopicNadir:                    prefix + "nadir",
		TopicNightBeforeStart:         prefix + "night/before/start",
		TopicNightBeforeEnd:           prefix + "night/before/end",
		TopicNightBeforeDuration:      prefix + "night/before/duration",
		TopicAstronomicalDawnStart:    prefix + "dawn/astronomical/start",
		TopicAstronomicalDawnEnd:      prefix + "dawn/astronomical/end",
		TopicAstronomicalDawnDuration: prefix + "dawn/astronomical/duration",
		TopicNauticalDawnStart:        prefix + "dawn/nautical/start",
		TopicNauticalDawnEnd:          prefix + "dawn/nautical/end",
		TopicNauticalDawnDuration:     prefix + "dawn/nautical/duration",
		TopicCivilDawnStart:           prefix + "dawn/civil/start",
		TopicCivilDawnEnd:             prefix + "dawn/civil/end",
		TopicCivilDawnDuration:        prefix + "dawn/civil/duration",
		TopicRiseStart:                prefix + "rise/start",
		TopicRiseEnd:                  prefix + "rise/end",
		TopicRiseDuration:             prefix + "rise/duration",
		TopicSolarNoon:                prefix + "solar-noon",
		TopicSetStart:                 prefix + "set/start",
		TopicSetEnd:                   prefix + "set/end",
		TopicSetDuration:              prefix + "set/duration",
		TopicCivilDuskStart:           prefix + "dusk/civil/start",
		TopicCivilDuskEnd:             prefix + "dusk/civil/end",
		TopicCivilDuskDuration:        prefix + "dusk/civil/duration",
		TopicNauticalDuskStart:        prefix + "dusk/nautical/start",
		TopicNauticalDuskEnd:          prefix + "dusk/nautical/end",
		TopicNauticalDuskDuration:     prefix + "dusk/nautical/duration",
		TopicAstronomicalDuskStart:    prefix + "dusk/astronomical/start",
		TopicAstronomicalDuskEnd:      prefix + "dusk/astronomical/end",
		TopicAstronomicalDuskDuration: prefix + "dusk/astronomical/duration",
		TopicNightAfterStart:          prefix + "night/after/start",
		TopicNightAfterEnd:            prefix + "night/after/end",
		TopicNightAfterDuration:       prefix + "night/after/duration",
	}
}
