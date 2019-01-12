package broadlink

type ConfigRM struct {
	IP  string `valid:"ip,required"`
	MAC string `valid:"mac,required"`
}

func (t TypeRM) Config() interface{} {
	return &ConfigRM{}
}
