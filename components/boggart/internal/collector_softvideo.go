package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
)

func (c *MetricsCollector) CollectSoftVideo() error {
	client := softvideo.NewClient(
		c.component.config.GetString(boggart.ConfigSoftVideoLogin),
		c.component.config.GetString(boggart.ConfigSoftVideoPassword))

	value, err := client.Balance()
	if err != nil {
		// FIXME: logging
		return err
	}

	metricSoftVideoBalance.Set(float64(value))

	return nil
}
