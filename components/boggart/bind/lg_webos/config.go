package lg_webos

type Config struct {
	Host string `valid:"host,required"`
	Key  string `valid:"required"`
}

func (t Type) Config() interface{} {
	return &Config{}
}
