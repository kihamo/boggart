package google_home

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/google/home"
)

type Bind struct {
	boggart.BindBase
	di.WorkersBind

	provider *home.Client
}
