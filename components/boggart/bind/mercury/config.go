package mercury

type Config struct {
	RS485 struct {
		Address string `valid:"required"`
		Timeout string
	} `valid:"required"`
	Address string `valid:"required"`
}

func (t Type) Config() interface{} {
	return &Config{}
}
