package owntracks

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config

	lat            *atomic.Float64
	lon            *atomic.Float64
	geoHash        *atomic.String
	conn           *atomic.String
	acc            *atomic.Int64
	alt            *atomic.Int64
	batt           *atomic.Float64
	vel            *atomic.Int64
	wayPointsCheck map[string]*atomic.BoolNull
}

func (b *Bind) SetStatusManager(getter boggart.BindStatusGetter, setter boggart.BindStatusSetter) {
	b.BindBase.SetStatusManager(getter, setter)

	b.UpdateStatus(boggart.BindStatusOnline)
}

func (b *Bind) validAccuracy(acc *int64) bool {
	if acc == nil {
		return false
	}

	value := *acc
	if value == 0 {
		return false
	}

	if b.config.MaxAccuracy > 0 && value > b.config.MaxAccuracy {
		return false
	}

	return true
}
