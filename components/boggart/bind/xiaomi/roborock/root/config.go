package root

const (
	DefaultDeviceIDFile      = "/mnt/data/miio/device.uid"
	DefaultRuntimeConfigFile = "/mnt/data/rockrobo/RoboController.cfg"
)

type Config struct {
	DeviceIDFile      string `mapstructure:"device_id_file" yaml:"device_id_file"`
	RuntimeConfigFile string `mapstructure:"runtime_config_file" yaml:"runtime_config_file"`
}

func (t Type) Config() interface{} {
	return &Config{
		DeviceIDFile:      DefaultDeviceIDFile,
		RuntimeConfigFile: DefaultRuntimeConfigFile,
	}
}
