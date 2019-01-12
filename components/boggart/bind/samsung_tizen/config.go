package samsung_tizen

type Config struct {
	Host string `valid:"host,required"`
}

func (t Type) Config() interface{} {
	return &Config{}
}
