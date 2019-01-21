package owntracks

type Config struct {
	Devices map[string]string
}

func (Type) Config() interface{} {
	return &Config{}
}
