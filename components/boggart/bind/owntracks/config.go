package owntracks

const (
	DefaultGeoFence = 100
)

type Point struct {
	Lat      float64
	Lon      float64
	GeoFence float64
}

type Config struct {
	User    string `valid:"required"`
	Device  string `valid:"required"`
	Regions map[string]Point
}

func (Type) Config() interface{} {
	return &Config{}
}
