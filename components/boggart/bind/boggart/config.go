package boggart

type Config struct {
	ApplicationName    string `valid:",required" mapstructure:"application_name" yaml:"application_name"`
	ApplicationVersion string `valid:",required" mapstructure:"application_version" yaml:"application_version"`
	ApplicationBuild   string `valid:",required" mapstructure:"application_build" yaml:"application_build"`
}

func (Type) Config() interface{} {
	return &Config{}
}
