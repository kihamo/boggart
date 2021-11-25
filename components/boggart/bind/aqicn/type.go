package aqicn

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/aqicn"
)

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{}
}

func (t Type) DashboardTemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		"icon": aqicn.Icon,
	}
}
