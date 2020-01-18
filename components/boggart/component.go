package boggart

import (
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/metrics"
)

type Component interface {
	shadow.Component
	metrics.HasMetrics

	ReloadConfig() (int, error)
	ReloadConfigByID(id string) error

	Bind(id string) BindItem
	BindItems() []BindItem
	RegisterBind(id string, bind Bind, t string, description string, tags []string, config interface{}) (BindItem, error)
	UnregisterBindByID(id string) error
}
