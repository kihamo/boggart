package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow"
)

type Component struct {
}

func (c *Component) GetName() string {
	return boggart.ComponentName
}

func (c *Component) GetVersion() string {
	return boggart.ComponentVersion
}

func (c *Component) GetDependencies() []shadow.Dependency {
	return []shadow.Dependency{}
}
