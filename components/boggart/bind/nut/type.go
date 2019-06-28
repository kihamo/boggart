package nut

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	bind := &Bind{
		config:    c.(*Config),
		variables: make(map[string]interface{}),
	}

	return bind, nil
}
