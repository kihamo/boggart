package hilink

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
)

const (
	DefaultLivenessInterval       = time.Minute
	DefaultLivenessTimeout        = time.Second * 5
	DefaultBalanceUpdaterInterval = time.Hour
	DefaultBalanceUpdaterTimeout  = time.Second * 30
	DefaultSMSCheckerInterval     = time.Minute
	DefaultSMSCheckerTimeout      = time.Second * 30
	DefaultSystemUpdaterInterval  = time.Minute
	DefaultSystemUpdaterTimeout   = time.Second * 30
)

type Config struct {
	Address                  boggart.URL `valid:",required"`
	Debug                    bool
	LivenessInterval         time.Duration `mapstructure:"liveness_interval" yaml:"liveness_interval"`
	LivenessTimeout          time.Duration `mapstructure:"liveness_timeout" yaml:"liveness_timeout"`
	BalanceUpdaterInterval   time.Duration `mapstructure:"balance_interval" yaml:"balance_interval"`
	BalanceUpdaterTimeout    time.Duration `mapstructure:"balance_timeout" yaml:"balance_timeout"`
	SMSCheckerInterval       time.Duration `mapstructure:"sms_checker_interval" yaml:"sms_checker_interval"`
	SMSCheckerTimeout        time.Duration `mapstructure:"sms_checker_timeout" yaml:"sms_checker_timeout"`
	SystemUpdaterInterval    time.Duration `mapstructure:"system_interval" yaml:"system_interval"`
	SystemUpdaterTimeout     time.Duration `mapstructure:"system_timeout" yaml:"system_timeout"`
	SMSCommandsEnabled       bool          `mapstructure:"sms_commands_enabled" yaml:"sms_commands_enabled"`
	SMSCommandsAllowedPhones []string      `mapstructure:"sms_commands_allowed_phones" yaml:"sms_commands_allowed_phones"`
}

func (t Type) Config() interface{} {
	return &Config{
		LivenessInterval:       DefaultLivenessInterval,
		LivenessTimeout:        DefaultLivenessTimeout,
		BalanceUpdaterInterval: DefaultBalanceUpdaterInterval,
		BalanceUpdaterTimeout:  DefaultBalanceUpdaterTimeout,
		SMSCheckerInterval:     DefaultSMSCheckerInterval,
		SMSCheckerTimeout:      DefaultSMSCheckerTimeout,
		SystemUpdaterInterval:  DefaultSystemUpdaterInterval,
		SystemUpdaterTimeout:   DefaultSystemUpdaterTimeout,
	}
}
