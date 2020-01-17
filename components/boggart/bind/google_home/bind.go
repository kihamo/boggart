package google_home

import (
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/google/home"
)

type Bind struct {
	di.MetaBind
	di.WorkersBind
	di.LoggerBind

	provider *home.Client
}
