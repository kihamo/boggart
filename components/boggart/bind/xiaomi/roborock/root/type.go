package root

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	bind := &Bind{
		config:             c.(*Config),
		cacheRuntimeConfig: make(map[string]string, 11),
		watchFiles:         make(map[string]func(string) error, 0),
	}

	return bind, nil
}
