package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/logger"
	"github.com/kihamo/shadow/components/metrics"
	"github.com/kihamo/shadow/components/workers"
)

type Component struct {
	config  config.Component
	workers workers.Component
	routes  []dashboard.Route
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

	return nil
}

func (c *Component) Run() (err error) {
	// запуск задач на снятие показаний

	return nil
}
