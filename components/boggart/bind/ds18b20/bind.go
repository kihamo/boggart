package ds18b20

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/yryz/ds18b20"
)

type Bind struct {
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	config *Config
}

func (b *Bind) Sensors() ([]string, error) {
	if len(b.config.Sensors) > 0 {
		return b.config.Sensors, nil
	}

	return ds18b20.Sensors()
}

func (b *Bind) Temperatures() (map[string]float64, error) {
	sensors, err := b.Sensors()
	if err != nil {
		return nil, err
	}

	values := make(map[string]float64, len(sensors))

	for _, sensor := range sensors {
		value, err := ds18b20.Temperature(sensor)
		if err != nil {
			return nil, err
		}

		values[sensor] = value
	}

	return values, nil
}
