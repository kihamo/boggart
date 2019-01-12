package nut

type Config struct {
	Host string `valid:"host,required"`
	UPS  string `valid:"required"`
}

func (t Type) Config() interface{} {
	return &Config{}
}
