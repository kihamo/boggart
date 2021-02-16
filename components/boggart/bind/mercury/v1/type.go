package v1

import (
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart"
	mercury "github.com/kihamo/boggart/providers/mercury/v1"
)

type Type struct {
	SerialNumberFunc func(address string) mercury.Option
	Device           uint8
}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{
		providerOnce: &atomic.Once{},
		tariffCount:  atomic.NewUint32Null(),
	}
}
