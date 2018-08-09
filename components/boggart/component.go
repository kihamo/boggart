package boggart

import (
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/metrics"
)

type Component interface {
	shadow.Component
	metrics.HasMetrics

	RS485() *rs485.Connection
}
