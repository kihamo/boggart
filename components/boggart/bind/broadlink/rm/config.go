package rm

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/types"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Host               string             `valid:",required"`
	MAC                types.HardwareAddr `valid:",required"`
	Model              string             `valid:"in(rm3mini|rm2proplus),required"`
	CaptureDuration    time.Duration      `mapstructure:"capture_interval" yaml:"capture_interval"`
	ConnectionTimeout  time.Duration      `mapstructure:"connection_timeout" yaml:"connection_timeout"`
	TopicCaptureState  mqtt.Topic         `mapstructure:"topic_capture_state" yaml:"topic_capture_state"`
	TopicCaptureSwitch mqtt.Topic         `mapstructure:"topic_capture_switch" yaml:"topic_capture_switch"`
	TopicCaptureResult mqtt.Topic         `mapstructure:"topic_capture_result" yaml:"topic_capture_result"`
	TopicIR            mqtt.Topic         `mapstructure:"topic_ir" yaml:"topic_ir"`
	TopicRF315         mqtt.Topic         `mapstructure:"topic_rf315" yaml:"topic_rf315"`
	TopicRF433         mqtt.Topic         `mapstructure:"topic_rf433" yaml:"topic_rf433"`
}

func (t Type) ConfigDefaults() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/remote-control/+/"

	probesConfig := di.ProbesConfigDefaults()
	probesConfig.ReadinessPeriod = time.Second * 30
	probesConfig.ReadinessTimeout = time.Second * 10

	return &Config{
		ProbesConfig:       probesConfig,
		LoggerConfig:       di.LoggerConfigDefaults(),
		CaptureDuration:    time.Second * 15,
		ConnectionTimeout:  time.Second,
		TopicCaptureState:  prefix + "capture",
		TopicCaptureSwitch: prefix + "capture/switch",
		TopicCaptureResult: prefix + "capture/result",
		TopicIR:            prefix + "ir",
		TopicRF315:         prefix + "rf315",
		TopicRF433:         prefix + "rf433",
	}
}
