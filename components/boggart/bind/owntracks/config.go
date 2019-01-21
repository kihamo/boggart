package owntracks

type Config struct {
	User   string `valid:"required"`
	Device string `valid:"required"`
}

func (Type) Config() interface{} {
	return &Config{}
}
