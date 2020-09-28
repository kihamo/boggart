package timelapse

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricTotalFiles = snitch.NewGauge(boggart.ComponentName+"_bind_timelapse_file_total", "Total files")
	metricTotalSize  = snitch.NewGauge(boggart.ComponentName+"_bind_timelapse_size_total_bytes", "Total sizes in bytes")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricTotalFiles.With("id", id).Describe(ch)
	metricTotalSize.With("id", id).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	id := b.Meta().ID()
	if id == "" {
		return
	}

	metricTotalFiles.With("id", id).Collect(ch)
	metricTotalSize.With("id", id).Collect(ch)
}
