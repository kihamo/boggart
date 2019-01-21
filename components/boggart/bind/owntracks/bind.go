package owntracks

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/atomic"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	user    string
	device  string
	regions map[string]Point

	lat     *atomic.Float64
	lon     *atomic.Float64
	geoHash *atomic.String
	conn    *atomic.String
	acc     *atomic.Int64
	alt     *atomic.Int64
	batt    *atomic.Float64
	vel     *atomic.Int64
	region  map[string]*atomic.Bool
}

func (b *Bind) SetStatusManager(getter boggart.BindStatusGetter, setter boggart.BindStatusSetter) {
	b.BindBase.SetStatusManager(getter, setter)

	b.UpdateStatus(boggart.BindStatusOnline)
}
