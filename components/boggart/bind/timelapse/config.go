package timelapse

import (
	"os"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	SaveDirectory  string      `valid:"required" mapstructure:"save_directory" yaml:"save_directory"`
	CaptureURL     boggart.URL `valid:"required" mapstructure:"capture_url" yaml:"capture_url"`
	FileNameFormat string      `mapstructure:"file_name_format" yaml:"file_name_format"`
}

func (Type) Config() interface{} {
	cacheDir, _ := os.UserCacheDir()
	if cacheDir == "" {
		cacheDir = os.TempDir()
	}

	if cacheDir != "" {
		cacheDirBind := cacheDir + string(os.PathSeparator) + boggart.ComponentName + "_timelapse"

		err := os.Mkdir(cacheDirBind, 0700)
		if err == nil || os.IsExist(err) {
			cacheDir = cacheDirBind
		}
	}

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Minute * 10,
			ReadinessTimeout: di.ProbesConfigLivenessDefaultTimeout,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		SaveDirectory:  cacheDir,
		FileNameFormat: "20060102_150405",
	}
}
