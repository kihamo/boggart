package broadlink

type ConfigSP3S struct {
	IP  string `valid:"ip,required"`
	MAC string `valid:"mac,required"`
}

func (t TypeSP3S) Config() interface{} {
	return &ConfigSP3S{}
}
