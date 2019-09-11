package internal

import (
	"github.com/kihamo/boggart/components/barcode"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/dashboard"
)

type Component struct {
	routes []dashboard.Route
}

func (c *Component) Name() string {
	return barcode.ComponentName
}

func (c *Component) Version() string {
	return barcode.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{}
}

func (c *Component) Run(_ shadow.Application, _ chan<- struct{}) error {
	return nil
}
