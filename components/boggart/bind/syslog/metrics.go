package syslog

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricHandledMessage = snitch.NewCounter(boggart.ComponentName+"_bind_syslog_handled_message_total", "Total handled messages")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricHandledMessage.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricHandledMessage.Collect(ch)
}
