package timelapse

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	SaveDirectory     string    `valid:"required" mapstructure:"save_directory" yaml:"save_directory"`
	CaptureURL        types.URL `valid:"required" mapstructure:"capture_url" yaml:"capture_url"`
	FileNameFormat    string    `mapstructure:"file_name_format" yaml:"file_name_format"`
	FileMode          uint32    `mapstructure:"file_mode" yaml:"file_mode"`
	SaveDirectoryMode uint32    `mapstructure:"save_directory_mode" yaml:"save_directory_mode"`
}

func (Type) Config() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute * 10

	return &Config{
		ProbesConfig:      probesConfig,
		LoggerConfig:      di.LoggerConfigDefaults(),
		FileNameFormat:    "20060102_150405",
		FileMode:          0664,
		SaveDirectoryMode: 0774,
	}
}
