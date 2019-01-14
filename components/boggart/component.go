package boggart

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/metrics"
)

type Component interface {
	shadow.Component
	metrics.HasMetrics

	ReloadConfig() error
}
