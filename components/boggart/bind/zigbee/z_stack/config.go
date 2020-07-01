package z_stack

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

type Config struct {
	di.ProbesConfig `mapstructure:",squash" yaml:",inline"`
	di.LoggerConfig `mapstructure:",squash" yaml:",inline"`

	ConnectionDSN                 string        `mapstructure:"connection_dsn" yaml:"connection_dsn" valid:"required"`
	DisableLED                    bool          `mapstructure:"disable_led" yaml:"disable_led"`
	PermitJoin                    bool          `mapstructure:"permit_join" yaml:"permit_join"`
	PermitJoinDuration            time.Duration `mapstructure:"permit_join_duration" yaml:"permit_join_duration"`
	TopicPermitJoin               mqtt.Topic    `mapstructure:"topic_permit_join" yaml:"topic_permit_join"`
	TopicVersionTransportRevision mqtt.Topic    `mapstructure:"topic_version_transport_revision" yaml:"topic_version_transport_revision"`
	TopicVersionProduct           mqtt.Topic    `mapstructure:"topic_version_product" yaml:"topic_version_product"`
	TopicVersionMajorRelease      mqtt.Topic    `mapstructure:"topic_version_major_release" yaml:"topic_version_major_release"`
	TopicVersionMinorRelease      mqtt.Topic    `mapstructure:"topic_version_minor_release" yaml:"topic_version_minor_release"`
	TopicVersionMainTrel          mqtt.Topic    `mapstructure:"topic_version_main_trel" yaml:"topic_version_main_trel"`
	TopicVersionHardwareRevision  mqtt.Topic    `mapstructure:"topic_version_hardware_revision" yaml:"topic_version_hardware_revision"`
	TopicVersionType              mqtt.Topic    `mapstructure:"topic_version_type" yaml:"topic_version_type"`
	TopicStatePermitJoin          mqtt.Topic    `mapstructure:"topic_state_permit_join" yaml:"topic_state_permit_join"`
	TopicStatePermitJoinDuration  mqtt.Topic    `mapstructure:"topic_state_permit_join_duration" yaml:"topic_state_permit_join_duration"`
	TopicLinkQuality              mqtt.Topic    `mapstructure:"topic_link_quality" yaml:"topic_link_quality"`
	TopicBatteryPercent           mqtt.Topic    `mapstructure:"topic_battery_percent" yaml:"topic_battery_percent"`
	TopicBatteryVoltage           mqtt.Topic    `mapstructure:"topic_battery_voltage" yaml:"topic_battery_voltage"`
	TopicOnOff                    mqtt.Topic    `mapstructure:"topic_on_off" yaml:"topic_on_off"`
	TopicClick                    mqtt.Topic    `mapstructure:"topic_click" yaml:"topic_click"`
}

func (Type) Config() interface{} {
	var (
		prefix       mqtt.Topic = boggart.ComponentName + "/zigbee/zstack/+/"
		prefixDevice            = prefix + "+/"
	)

	return &Config{
		ProbesConfig: di.ProbesConfig{
			ReadinessPeriod:  time.Second * 30,
			ReadinessTimeout: time.Second * 5,
			LivenessPeriod:   time.Minute * 10,
			LivenessTimeout:  time.Second * 5,
		},
		LoggerConfig: di.LoggerConfig{
			BufferedRecordsLimit: di.LoggerDefaultBufferedRecordsLimit,
			BufferedRecordsLevel: di.LoggerDefaultBufferedRecordsLevel,
		},
		DisableLED:                    false,
		PermitJoin:                    false,
		PermitJoinDuration:            255 * time.Second,
		TopicPermitJoin:               prefix + "permit-join",
		TopicVersionTransportRevision: prefix + "version/transport-revision",
		TopicVersionProduct:           prefix + "version/product",
		TopicVersionMajorRelease:      prefix + "version/major-release",
		TopicVersionMinorRelease:      prefix + "version/minor-release",
		TopicVersionMainTrel:          prefix + "version/main-trel",
		TopicVersionHardwareRevision:  prefix + "version/hardware-revision",
		TopicVersionType:              prefix + "version/type",
		TopicStatePermitJoin:          prefix + "state/permit-join",
		TopicStatePermitJoinDuration:  prefix + "state/permit-join-duration",
		TopicLinkQuality:              prefixDevice + "link-quality",
		TopicBatteryPercent:           prefixDevice + "battery/percent",
		TopicBatteryVoltage:           prefixDevice + "battery/voltage",
		TopicOnOff:                    prefixDevice + "on-off",
		TopicClick:                    prefixDevice + "click",
	}
}
