package gpio

type Config struct {
	Pin  uint64 `valid:"required"`
	Mode string `valid:"in(in|out)"`
}

func (t Type) Config() interface{} {
	return &Config{}
}
