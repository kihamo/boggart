package mikrotik

type Config struct {
	Address      string `valid:"url,required"`
	SyslogClient string `valid:"host"`
}

func (t Type) Config() interface{} {
	return &Config{}
}
