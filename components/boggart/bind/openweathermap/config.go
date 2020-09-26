package openweathermap

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/openweathermap"
)

type Config struct {
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Debug                 bool
	UpdaterInterval       time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	APIKey                string        `mapstructure:"api_key" yaml:"api_key" valid:",required"`
	Units                 string
	CityID                uint64 `mapstructure:"city_id" yaml:"city_id"`
	CityName              string `mapstructure:"city_name" yaml:"city_name"`
	Latitude              float64
	Longitude             float64
	Zip                   string `mapstructure:"zip" yaml:"zip"`
	Price                 int
	TopicCurrentTemp      mqtt.Topic `mapstructure:"topic_current_temp" yaml:"topic_current_temp"`
	TopicDailyTempMin     mqtt.Topic `mapstructure:"topic_daily_temp_min" yaml:"topic_daily_temp_min"`
	TopicDailyTempMax     mqtt.Topic `mapstructure:"topic_daily_temp_max" yaml:"topic_daily_temp_max"`
	TopicDailyTempDay     mqtt.Topic `mapstructure:"topic_daily_temp_day" yaml:"topic_daily_temp_day"`
	TopicDailyTempNight   mqtt.Topic `mapstructure:"topic_daily_temp_night" yaml:"topic_daily_temp_night"`
	TopicDailyTempMorning mqtt.Topic `mapstructure:"topic_daily_temp_morning" yaml:"topic_daily_temp_morning"`
	TopicDailyWindSpeed   mqtt.Topic `mapstructure:"topic_daily_wind_speed" yaml:"topic_daily_wind_speed"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/service/openweathermap/+/"

	return &Config{
		LoggerConfig:          di.LoggerConfigDefaults(),
		Units:                 "metric",
		Price:                 openweathermap.PriceFree,
		UpdaterInterval:       time.Minute * 15,
		TopicCurrentTemp:      prefix + "current",
		TopicDailyTempMin:     prefix + "daily/+/min",
		TopicDailyTempMax:     prefix + "daily/+/max",
		TopicDailyTempDay:     prefix + "daily/+/day",
		TopicDailyTempNight:   prefix + "daily/+/night",
		TopicDailyTempMorning: prefix + "daily/+/morning",
		TopicDailyWindSpeed:   prefix + "daily/+/wind/speed",
	}
}
