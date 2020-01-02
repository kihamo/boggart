package scale

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	MAC                    boggart.HardwareAddr `valid:",required"`
	UpdaterInterval        time.Duration        `mapstructure:"updater_interval" yaml:"updater_interval"`
	CaptureDuration        time.Duration        `mapstructure:"capture_interval" yaml:"capture_interval"`
	TopicDatetime          mqtt.Topic           `mapstructure:"topic_datetime" yaml:"topic_datetime"`
	TopicWeight            mqtt.Topic           `mapstructure:"topic_weight" yaml:"topic_weight"`
	TopicImpedance         mqtt.Topic           `mapstructure:"topic_impedance" yaml:"topic_impedance"`
	TopicProfile           mqtt.Topic           `mapstructure:"topic_profile" yaml:"topic_profile"`
	TopicProfileSet        mqtt.Topic           `mapstructure:"topic_profile_set" yaml:"topic_profile_set"`
	TopicBMR               mqtt.Topic           `mapstructure:"topic_bmr" yaml:"topic_bmr"`
	TopicBMI               mqtt.Topic           `mapstructure:"topic_bmi" yaml:"topic_bmi"`
	TopicFatPercentage     mqtt.Topic           `mapstructure:"topic_fat_percentage" yaml:"topic_fat_percentage"`
	TopicWaterPercentage   mqtt.Topic           `mapstructure:"topic_water_percentage" yaml:"topic_water_percentage"`
	TopicIdealWeight       mqtt.Topic           `mapstructure:"topic_ideal_weight" yaml:"topic_ideal_weight"`
	TopicLBMCoefficient    mqtt.Topic           `mapstructure:"topic_lbm_coefficient" yaml:"topic_lbm_coefficient"`
	TopicBoneMass          mqtt.Topic           `mapstructure:"topic_bone_mass" yaml:"topic_bone_mass"`
	TopicMuscleMass        mqtt.Topic           `mapstructure:"topic_muscle_mass" yaml:"topic_muscle_mass"`
	TopicVisceralFat       mqtt.Topic           `mapstructure:"topic_visceral_fat" yaml:"topic_visceral_fat"`
	TopicFatMassToIdeal    mqtt.Topic           `mapstructure:"topic_fat_mass_to_ideal" yaml:"topic_fat_mass_to_ideal"`
	TopicProteinPercentage mqtt.Topic           `mapstructure:"topic_protein_percentage" yaml:"topic_protein_percentage"`
	TopicBodyType          mqtt.Topic           `mapstructure:"topic_body_type" yaml:"topic_body_type"`
	TopicMetabolicAge      mqtt.Topic           `mapstructure:"topic_metabolic_age" yaml:"topic_metabolic_age"`
}

func (Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/xiaomi/scale/+/"

	return &Config{
		UpdaterInterval:        time.Minute,
		CaptureDuration:        time.Second * 10,
		TopicDatetime:          prefix + "datetime",
		TopicWeight:            prefix + "weight",
		TopicImpedance:         prefix + "impedance",
		TopicProfile:           prefix + "profile",
		TopicProfileSet:        prefix + "profile/set",
		TopicBMR:               prefix + "bmr",
		TopicBMI:               prefix + "bmi",
		TopicFatPercentage:     prefix + "fat-percentage",
		TopicWaterPercentage:   prefix + "water-percentage",
		TopicIdealWeight:       prefix + "ideal-weight",
		TopicLBMCoefficient:    prefix + "lbm-coefficient",
		TopicBoneMass:          prefix + "bone-mass",
		TopicMuscleMass:        prefix + "muscle-mass",
		TopicVisceralFat:       prefix + "visceral-fat",
		TopicFatMassToIdeal:    prefix + "fat-mass-to-ideal",
		TopicProteinPercentage: prefix + "protein-percentage",
		TopicBodyType:          prefix + "body-type",
		TopicMetabolicAge:      prefix + "metabolic-age",
	}
}
