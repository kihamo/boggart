package owntracks

import (
	"github.com/kihamo/boggart/components/boggart"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	user   string
	device string
}

func (b *Bind) SetStatusManager(getter boggart.BindStatusGetter, setter boggart.BindStatusSetter) {
	b.BindBase.SetStatusManager(getter, setter)

	b.UpdateStatus(boggart.BindStatusOnline)
}
