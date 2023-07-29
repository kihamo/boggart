package timelapse

import (
	"os"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	UpdaterInterval       time.Duration  `mapstructure:"updater_interval" yaml:"updater_interval"`
	SaveDirectory         string         `valid:"required" mapstructure:"save_directory" yaml:"save_directory"`
	CaptureURL            types.URL      `valid:"required" mapstructure:"capture_url" yaml:"capture_url"`
	FileNameFormat        string         `mapstructure:"file_name_format" yaml:"file_name_format"`
	FilePerm              types.FileMode `mapstructure:"file_perm" yaml:"file_perm"`
	DirectoryPerm         types.FileMode `mapstructure:"directory_perm" yaml:"directory_perm"`
	FilesOnPage           int            `mapstructure:"files_on_page" yaml:"files_on_page"`
	EnableMigrationV1ToV2 bool           `mapstructure:"enable_migration_v1_to_v2" yaml:"enable_migration_v1_to_v2"`
	TopicCapture          mqtt.Topic     `mapstructure:"topic_capture" yaml:"topic_capture"`
}

func (Type) ConfigDefaults() interface{} {
	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Minute * 10

	return &Config{
		ProbesConfig:    probesConfig,
		LoggerConfig:    di.LoggerConfigDefaults(),
		UpdaterInterval: time.Hour,
		FileNameFormat:  "20060102_150405",
		FilePerm: types.FileMode{
			FileMode: os.FileMode(0664),
		},
		DirectoryPerm: types.FileMode{
			FileMode: os.FileMode(0774),
		},
		FilesOnPage:  9,
		TopicCapture: boggart.ComponentName + "/timelapse/+/capture",
	}
}
