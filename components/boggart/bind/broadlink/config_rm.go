package broadlink

type ConfigRM struct {
	IP              string `valid:"ip,required"`
	MAC             string `valid:"mac,required"`
	Model           string `valid:"in(rm3mini|rm2proplus),required"`
	CaptureDuration string
}

func (t TypeRM) Config() interface{} {
	return &ConfigRM{
		CaptureDuration: RMCaptureDuration.String(),
	}
}
