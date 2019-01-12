package softvideo

type Config struct {
	Login           string `valid:"required"`
	Password        string `valid:"required"`
	UpdaterInterval string `mapstructure:"updater_interval"`
}

func (Type) Config() interface{} {
	return &Config{
		UpdaterInterval: DefaultUpdaterInterval.String(),
	}
}
