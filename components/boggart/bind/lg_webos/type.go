package lg_webos

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config:       config,
		quitMonitors: make(chan struct{}),
	}

	return bind, nil
}
