package scale

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/types"
)

type Profile struct {
	Sex      bool      `json:"sex,omitempty" yaml:"sex,omitempty"`
	Editable bool      `json:"-" yaml:"-"`
	Name     string    `json:"name,omitempty" yaml:"name,omitempty"`
	Height   uint64    `json:"height,omitempty" yaml:"height,omitempty"`
	Age      uint64    `json:"age,omitempty" yaml:"age,omitempty"`
	Birthday time.Time `json:"birthday,omitempty" yaml:"birthday,omitempty"`
}

func (p Profile) GetAge() (age uint64) {
	if p.Age > 0 {
		age = p.Age
	} else if !p.Birthday.IsZero() {
		now := time.Now()

		age = uint64(now.Year() - p.Birthday.Year())

		diff := time.Date(now.Year(), p.Birthday.Month(), p.Birthday.Day(), p.Birthday.Hour(), p.Birthday.Minute(), p.Birthday.Second(), p.Birthday.Nanosecond(), p.Birthday.Location())
		if diff.After(now) {
			age--
		}
	}

	return age
}

var profileGuest = &Profile{
	Name:     "guest",
	Editable: true,
}

type Config struct {
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	MAC                    types.HardwareAddr  `valid:",required"`
	Profiles               map[string]*Profile `valid:",required"`
	UpdaterInterval        time.Duration       `mapstructure:"updater_interval" yaml:"updater_interval"`
	CaptureDuration        time.Duration       `mapstructure:"capture_interval" yaml:"capture_interval"`
	IgnoreEmptyImpedance   bool                `mapstructure:"ignore_empty_impedance,omitempty" yaml:"ignore_empty_impedance"`
	TopicProfile           mqtt.Topic          `mapstructure:"topic_profile" yaml:"topic_profile"`
	TopicProfileSettings   mqtt.Topic          `mapstructure:"topic_profile_settings" yaml:"topic_profile_settings"`
	TopicProfileActivate   mqtt.Topic          `mapstructure:"topic_profile_activate" yaml:"topic_profile_activate"`
	TopicDatetime          mqtt.Topic          `mapstructure:"topic_datetime" yaml:"topic_datetime"`
	TopicWeight            mqtt.Topic          `mapstructure:"topic_weight" yaml:"topic_weight"`
	TopicImpedance         mqtt.Topic          `mapstructure:"topic_impedance" yaml:"topic_impedance"`
	TopicBMR               mqtt.Topic          `mapstructure:"topic_bmr" yaml:"topic_bmr"`
	TopicBMI               mqtt.Topic          `mapstructure:"topic_bmi" yaml:"topic_bmi"`
	TopicFatPercentage     mqtt.Topic          `mapstructure:"topic_fat_percentage" yaml:"topic_fat_percentage"`
	TopicWaterPercentage   mqtt.Topic          `mapstructure:"topic_water_percentage" yaml:"topic_water_percentage"`
	TopicIdealWeight       mqtt.Topic          `mapstructure:"topic_ideal_weight" yaml:"topic_ideal_weight"`
	TopicLBMCoefficient    mqtt.Topic          `mapstructure:"topic_lbm_coefficient" yaml:"topic_lbm_coefficient"`
	TopicBoneMass          mqtt.Topic          `mapstructure:"topic_bone_mass" yaml:"topic_bone_mass"`
	TopicMuscleMass        mqtt.Topic          `mapstructure:"topic_muscle_mass" yaml:"topic_muscle_mass"`
	TopicVisceralFat       mqtt.Topic          `mapstructure:"topic_visceral_fat" yaml:"topic_visceral_fat"`
	TopicFatMassToIdeal    mqtt.Topic          `mapstructure:"topic_fat_mass_to_ideal" yaml:"topic_fat_mass_to_ideal"`
	TopicProteinPercentage mqtt.Topic          `mapstructure:"topic_protein_percentage" yaml:"topic_protein_percentage"`
	TopicBodyType          mqtt.Topic          `mapstructure:"topic_body_type" yaml:"topic_body_type"`
	TopicMetabolicAge      mqtt.Topic          `mapstructure:"topic_metabolic_age" yaml:"topic_metabolic_age"`
}

func (Type) ConfigDefaults() interface{} {
	var (
		prefix        mqtt.Topic = boggart.ComponentName + "/xiaomi/scale/+/"
		prefixProfile            = prefix + "+/"
	)

	return &Config{
		LoggerConfig:           di.LoggerConfigDefaults(),
		UpdaterInterval:        time.Minute,
		CaptureDuration:        time.Second * 10,
		IgnoreEmptyImpedance:   true,
		TopicProfile:           prefix + "profile",
		TopicProfileActivate:   prefixProfile + "activate",
		TopicProfileSettings:   prefix + mqtt.Topic(profileGuest.Name) + "/settings",
		TopicDatetime:          prefixProfile + "datetime",
		TopicWeight:            prefixProfile + "weight",
		TopicImpedance:         prefixProfile + "impedance",
		TopicBMR:               prefixProfile + "bmr",
		TopicBMI:               prefixProfile + "bmi",
		TopicFatPercentage:     prefixProfile + "fat-percentage",
		TopicWaterPercentage:   prefixProfile + "water-percentage",
		TopicIdealWeight:       prefixProfile + "ideal-weight",
		TopicLBMCoefficient:    prefixProfile + "lbm-coefficient",
		TopicBoneMass:          prefixProfile + "bone-mass",
		TopicMuscleMass:        prefixProfile + "muscle-mass",
		TopicVisceralFat:       prefixProfile + "visceral-fat",
		TopicFatMassToIdeal:    prefixProfile + "fat-mass-to-ideal",
		TopicProteinPercentage: prefixProfile + "protein-percentage",
		TopicBodyType:          prefixProfile + "body-type",
		TopicMetabolicAge:      prefixProfile + "metabolic-age",
	}
}
