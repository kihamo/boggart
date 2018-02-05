package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/softvideo"
	"github.com/kihamo/snitch"
)

const (
	MetricSoftVideoBalance = boggart.ComponentName + "_softvideo_balance_rubles_total"
)

var (
	metricSoftVideoBalance = snitch.NewGauge(MetricSoftVideoBalance, "SoftVideo balance in rubles")
)

func (c *MetricsCollector) UpdaterSoftVideo() error {
	client := softvideo.NewClient(
		c.component.config.String(boggart.ConfigSoftVideoLogin),
		c.component.config.String(boggart.ConfigSoftVideoPassword))

	value, err := client.Balance()
	if err != nil {
		// FIXME: logging
		return err
	}

	metricSoftVideoBalance.Set(float64(value))

	return nil
}

func (c *MetricsCollector) DescribeSoftVideo(ch chan<- *snitch.Description) {
	metricSoftVideoBalance.Describe(ch)
}

func (c *MetricsCollector) CollectSoftVideo(ch chan<- snitch.Metric) {
	metricSoftVideoBalance.Collect(ch)
}
