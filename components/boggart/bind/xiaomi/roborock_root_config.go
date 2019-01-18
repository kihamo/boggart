package xiaomi

const (
	DefaultDeviceIDFile      = "/mnt/data/miio/device.uid"
	DefaultRuntimeConfigFile = "/mnt/data/rockrobo/RoboController.cfg"
)

type RoborockRootConfig struct {
	DeviceIDFile      string `mapstructure:"device_id_file" yaml:"device_id_file"`
	RuntimeConfigFile string `mapstructure:"runtime_config_file" yaml:"runtime_config_file"`
}

func (t RoborockRootType) Config() interface{} {
	return &RoborockRootConfig{
		DeviceIDFile:      DefaultDeviceIDFile,
		RuntimeConfigFile: DefaultRuntimeConfigFile,
	}
}
