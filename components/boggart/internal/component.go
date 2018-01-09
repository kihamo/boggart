package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

type Component struct {
	config config.Component
	routes []dashboard.Route
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
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.config = a.GetComponent(config.ComponentName).(config.Component)

	return nil
}
