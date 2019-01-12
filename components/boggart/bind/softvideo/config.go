package softvideo

type Config struct {
	Login    string `valid:"required"`
	Password string `valid:"required"`
}

func (Type) Config() interface{} {
	return &Config{}
}
