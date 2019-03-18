package homie

const (
	DefaultBaseTopic = "homie"
)

type Config struct {
	DeviceID  string `valid:"required" mapstructure:"device_id" yaml:"device_id"`
	BaseTopic string `mapstructure:"base_topic" yaml:"base_topic"`
}

func (t Type) Config() interface{} {
	return &Config{
		BaseTopic: DefaultBaseTopic,
	}
}
