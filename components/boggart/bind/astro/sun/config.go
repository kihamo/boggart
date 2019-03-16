package sun

type Config struct {
	Lat  float64 `valid:"required"`
	Lon  float64 `valid:"required"`
	Name string  `valid:"required"`
}

func (t Type) Config() interface{} {
	return &Config{}
}
