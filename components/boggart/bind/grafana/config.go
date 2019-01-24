package grafana

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Config struct {
	Address    boggart.URL `valid:",required"`
	Name       string      `valid:",required"`
	ApiKey     string      `mapstructure:"api_key" yaml:"api_key"`
	Username   string
	Password   string
	Dashboards []int64
}

func (t Type) Config() interface{} {
	return &Config{}
}
