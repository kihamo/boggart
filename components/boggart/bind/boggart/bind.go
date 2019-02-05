package boggart

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config *Config
}

func (b *Bind) SetStatusManager(getter boggart.BindStatusGetter, setter boggart.BindStatusSetter) {
	b.BindBase.SetStatusManager(getter, setter)

	b.UpdateStatus(boggart.BindStatusOnline)
}
