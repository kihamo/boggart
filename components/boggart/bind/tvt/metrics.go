package tvt

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricStorageUsage     = snitch.NewGauge(boggart.ComponentName+"_bind_tvt_storage_usage_bytes", "TVT storage usage in bytes")
	metricStorageAvailable = snitch.NewGauge(boggart.ComponentName+"_bind_tvt_storage_available_bytes", "TVT storage available in bytes")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricStorageUsage.With("serial_number", sn).Describe(ch)
	metricStorageAvailable.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()
	if sn == "" {
		return
	}

	metricStorageUsage.With("serial_number", sn).Collect(ch)
	metricStorageAvailable.With("serial_number", sn).Collect(ch)
}
