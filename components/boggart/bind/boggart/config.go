package boggart

type Config struct {
	Build string
}

func (Type) Config() interface{} {
	return &Config{}
}
