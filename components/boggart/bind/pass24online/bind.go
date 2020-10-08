package pass24online

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/pass24online"
)

type Bind struct {
	di.LoggerBind
	di.WidgetBind

	provider *pass24online.Client
}
