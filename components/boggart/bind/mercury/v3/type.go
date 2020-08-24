package v3

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	return &Bind{
		config:       c.(*Config),
		providerOnce: &atomic.Once{},
	}, nil
}
