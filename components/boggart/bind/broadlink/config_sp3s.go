package broadlink

type ConfigSP3S struct {
	IP              string `valid:"ip,required"`
	MAC             string `valid:"mac,required"`
	UpdaterInterval string `mapstructure:"updater_interval"`
}

func (t TypeSP3S) Config() interface{} {
	return &ConfigSP3S{
		UpdaterInterval: SP3SDefaultUpdateInterval.String(),
	}
}
