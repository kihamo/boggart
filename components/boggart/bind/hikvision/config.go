package hikvision

type Config struct {
	Address string `valid:"url,required"`
}

func (t Type) Config() interface{} {
	return &Config{}
}
