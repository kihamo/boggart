package sun

import (
	"time"
)

type Config struct {
	Lat       float64 `valid:"required"`
	Lon       float64 `valid:"required"`
	UTCOffset float64 `mapstructure:"utc_offset" yaml:"utc_offset"`
	Name      string  `valid:"required"`
}

func (t Type) Config() interface{} {
	_, offset := time.Now().Zone()

	return &Config{
		UTCOffset: float64(offset / (60 * 60)),
	}
}
