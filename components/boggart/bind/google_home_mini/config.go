package google_home_mini

type Config struct {
	Host string `valid:"host,required"`
}

func (t Type) Config() interface{} {
	return &Config{}
}
