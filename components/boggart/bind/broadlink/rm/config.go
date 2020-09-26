package rm

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/types"
)

type ConfigRM struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	Host                 string             `valid:",required"`
	MAC                  types.HardwareAddr `valid:",required"`
	Model                string             `valid:"in(rm3mini|rm2proplus),required"`
	CaptureDuration      time.Duration      `mapstructure:"capture_interval" yaml:"capture_interval"`
	ConnectionTimeout    time.Duration      `mapstructure:"connection_timeout" yaml:"connection_timeout"`
	TopicCapture         mqtt.Topic         `mapstructure:"topic_capture" yaml:"topic_capture"`
	TopicCaptureState    mqtt.Topic         `mapstructure:"topic_capture_state" yaml:"topic_capture_state"`
	TopicIR              mqtt.Topic         `mapstructure:"topic_ir" yaml:"topic_ir"`
	TopicIRCount         mqtt.Topic         `mapstructure:"topic_ir_count" yaml:"topic_ir_count"`
	TopicIRCapture       mqtt.Topic         `mapstructure:"topic_ir_capture" yaml:"topic_ir_capture"`
	TopicRF315mhz        mqtt.Topic         `mapstructure:"topic_rf315mhz" yaml:"topic_rf315mhz"`
	TopicRF315mhzCount   mqtt.Topic         `mapstructure:"topic_rf315mhz_count" yaml:"topic_rf315mhz_count"`
	TopicRF315mhzCapture mqtt.Topic         `mapstructure:"topic_rf315mhz_capture" yaml:"topic_rf315mhz_capture"`
	TopicRF433mhz        mqtt.Topic         `mapstructure:"topic_rf433mhz" yaml:"topic_rf433mhz"`
	TopicRF433mhzCount   mqtt.Topic         `mapstructure:"topic_rf433mhz_count" yaml:"topic_rf433mhz_count"`
	TopicRF433mhzCapture mqtt.Topic         `mapstructure:"topic_rf433mhz_capture" yaml:"topic_rf433mhz_capture"`
}

func (t Type) Config() interface{} {
	var prefix mqtt.Topic = boggart.ComponentName + "/remote-control/+/"

	return &ConfigRM{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Second * 30,
			ReadinessTimeout: time.Second * 10,
		},
		LoggerConfig:         di.LoggerConfigDefaults(),
		CaptureDuration:      time.Second * 15,
		ConnectionTimeout:    time.Second,
		TopicCapture:         prefix + "capture",
		TopicCaptureState:    prefix + "capture/state",
		TopicIR:              prefix + "ir",
		TopicIRCount:         prefix + "ir/count",
		TopicIRCapture:       prefix + "capture/ir",
		TopicRF315mhz:        prefix + "rf315mhz",
		TopicRF315mhzCount:   prefix + "rf315mhz/count",
		TopicRF315mhzCapture: prefix + "capture/rf315mhz",
		TopicRF433mhz:        prefix + "rf433mhz",
		TopicRF433mhzCount:   prefix + "rf433mhz/count",
		TopicRF433mhzCapture: prefix + "capture/rf433mhz",
	}
}
