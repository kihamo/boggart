package internal

import (
	"github.com/kihamo/boggart/components/storage"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

type Component struct {
	routes []dashboard.Route
}

func (c *Component) Name() string {
	return storage.ComponentName
}

func (c *Component) Version() string {
	return storage.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	return nil
}
