package led_wifi

type Config struct {
	Address string `valid:"host,required"`
}

func (b Bind) Config() interface{} {
	return &Config{}
}
