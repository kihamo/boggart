package pulsar

import (
	"time"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	loc, err := time.LoadLocation(config.Location)
	if err != nil {
		return nil, err
	}

	bind := &Bind{
		config:         config,
		location:       loc,
		connectionOnce: &atomic.Once{},
	}

	return bind, nil
}
