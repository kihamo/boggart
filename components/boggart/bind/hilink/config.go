package hilink

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

const (
	DefaultDebug              = false
	DefaultBalanceRegexp      = `OCTATOK (?P<balance>\d+\.\d{2})\sp\..*?`
	DefaultBalanceUSSD        = "*105#"
	DefaultLivenessInterval   = time.Minute
	DefaultLivenessTimeout    = time.Second * 5
	DefaultUpdaterInterval    = time.Hour
	DefaultUpdaterTimeout     = time.Second * 30
	DefaultSMSCheckerInterval = time.Minute
	DefaultSMSCheckerTimeout  = time.Second * 30
)

type Config struct {
	Address            boggart.URL `valid:",required"`
	BalanceRegexp      string      `mapstructure:"balance_regexp" yaml:"balance_regexp"`
	BalanceUSSD        string      `mapstructure:"balance_ussd" yaml:"balance_ussd"`
	Debug              bool
	LivenessInterval   time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout    time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	UpdaterInterval    time.Duration `mapstructure:"updater_interval" yaml:"updater_interval"`
	UpdaterTimeout     time.Duration `mapstructure:"updater_timeout" yaml:"updater_timeout"`
	SMSCheckerInterval time.Duration `mapstructure:"sms_checker_interval" yaml:"sms_checker_interval"`
	SMSCheckerTimeout  time.Duration `mapstructure:"sms_checker_timeout" yaml:"sms_checker_timeout"`
}

func (t Type) Config() interface{} {
	return &Config{
		BalanceRegexp:      DefaultBalanceRegexp,
		BalanceUSSD:        DefaultBalanceUSSD,
		Debug:              DefaultDebug,
		LivenessInterval:   DefaultLivenessInterval,
		LivenessTimeout:    DefaultLivenessTimeout,
		UpdaterInterval:    DefaultUpdaterInterval,
		UpdaterTimeout:     DefaultUpdaterTimeout,
		SMSCheckerInterval: DefaultSMSCheckerInterval,
		SMSCheckerTimeout:  DefaultSMSCheckerTimeout,
	}
}
