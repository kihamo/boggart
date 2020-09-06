package openweathermap

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
)

type Config struct {
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Debug           bool
	UpdaterInterval time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	APIKey          string        `mapstructure:"api_key" yaml:"api_key" valid:",required"`
	Units           string
	CityID          uint64 `mapstructure:"city_id" yaml:"city_id"`
	CityName        string `mapstructure:"city_name" yaml:"city_name"`
	Latitude        float64
	Longitude       float64
	Zip             string `mapstructure:"zip" yaml:"zip"`
}

func (t Type) Config() interface{} {
	return &Config{
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		Units:           "metric",
		UpdaterInterval: time.Minute * 10,
	}
}
