package syslog

import (
	"github.com/kihamo/snitch"
)

var (
	metricHandledMessage = snitch.NewCounter("handled_message_total", "Total handled messages")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	metricHandledMessage.Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	metricHandledMessage.Collect(ch)
}
