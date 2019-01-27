package network

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type TypePing struct{}

func (t TypePing) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*ConfigPing)

	bind := &BindPing{
		hostname:        config.Hostname,
		retry:           config.Retry,
		timeout:         config.Timeout,
		updaterInterval: config.UpdaterInterval,
		online:          atomic.NewBoolNull(),
		latency:         atomic.NewUint32Null(),
	}

	return bind, nil
}
