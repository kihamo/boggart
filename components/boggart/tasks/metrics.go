package tasks

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricWorkers      = snitch.NewGauge(boggart.ComponentName+"_tasks_manager_workers_total", "Workers total")
	metricHandleStatus = snitch.NewGauge(boggart.ComponentName+"_tasks_manager_handle_success_total", "Handle status total")
)

func (m *Manager) Describe(ch chan<- *snitch.Description) {
	metricWorkers.Describe(ch)
	metricHandleStatus.Describe(ch)
}

func (m *Manager) Collect(ch chan<- snitch.Metric) {
	metricWorkers.Collect(ch)
	metricHandleStatus.Collect(ch)
}
