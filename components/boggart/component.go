package boggart

import (
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/metrics"
	"periph.io/x/periph/conn/onewire"
)

type Component interface {
	shadow.Component
	metrics.HasMetrics

	RS485() *rs485.Connection
	OneWire() onewire.BusCloser
}
