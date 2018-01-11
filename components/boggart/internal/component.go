package internal

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers/task"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
)

type Component struct {
	config    config.Component
	workers   workers.Component
	routes    []dashboard.Route
	collector *MetricsCollector
}

func (c *Component) GetName() string {
	return boggart.ComponentName
}

func (c *Component) GetVersion() string {
	return boggart.ComponentVersion
}

func (c *Component) GetDependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: dashboard.ComponentName,
		},
		{
			Name: logger.ComponentName,
		},
		{
			Name: metrics.ComponentName,
		},
		{
			Name:     workers.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.workers = a.GetComponent(workers.ComponentName).(workers.Component)
	c.collector = NewMetricsCollector(c)

	return nil
}

func (c *Component) Run() error {
	taskSoftVideo := task.NewFunctionTask(c.taskSoftVideo)
	taskSoftVideo.SetRepeats(-1)
	taskSoftVideo.SetRepeatInterval(c.config.GetDuration(boggart.ConfigSoftVideoRepeatInterval))
	taskSoftVideo.SetName(c.GetName() + "-softvideo-updater")
	c.workers.AddTask(taskSoftVideo)

	taskPulsar := task.NewFunctionTask(c.taskPulsar)
	taskPulsar.SetRepeats(-1)
	taskPulsar.SetRepeatInterval(c.config.GetDuration(boggart.ConfigPulsarRepeatInterval))
	taskPulsar.SetName(c.GetName() + "-pulsar-updater")
	c.workers.AddTask(taskPulsar)

	taskMikrotik := task.NewFunctionTask(c.taskMikrotik)
	taskMikrotik.SetRepeats(-1)
	taskMikrotik.SetRepeatInterval(c.config.GetDuration(boggart.ConfigMikrotikRepeatInterval))
	taskMikrotik.SetName(c.GetName() + "-mikrotik-updater")
	c.workers.AddTask(taskMikrotik)

	return nil
}

func (c *Component) taskPulsar(context.Context) (interface{}, error) {
	if c.config.GetBool(boggart.ConfigPulsarEnabled) {
		return nil, c.collector.CollectPulsar()
	}

	return nil, nil
}

func (c *Component) taskSoftVideo(context.Context) (interface{}, error) {
	if c.config.GetBool(boggart.ConfigSoftVideoEnabled) {
		return nil, c.collector.CollectSoftVideo()
	}

	return nil, nil
}

func (c *Component) taskMikrotik(context.Context) (interface{}, error) {
	if c.config.GetBool(boggart.ConfigMikrotikEnabled) {
		return nil, c.collector.CollectMikrotik()
	}

	return nil, nil
}
