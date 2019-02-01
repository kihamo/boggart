package fcm

type Config struct {
	Tokens      []string `valid:"required"`
	Credentials string   `valid:"required"`
}

func (Type) Config() interface{} {
	return &Config{}
}
