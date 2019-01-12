package ds18b20

type Config struct {
	Address string `valid:"required"`
}

func (t Type) Config() interface{} {
	return &Config{}
}
