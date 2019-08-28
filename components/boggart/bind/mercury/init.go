package mercury

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/mercury/v1"
	"github.com/kihamo/boggart/components/boggart/bind/mercury/v3"
)

func init() {
	version1 := v1.Type{}
	boggart.RegisterBindType("mercury:200", version1)
	boggart.RegisterBindType("mercury:201", version1)
	boggart.RegisterBindType("mercury:203", version1)
	boggart.RegisterBindType("mercury:206", version1)

	version3 := v3.Type{}
	boggart.RegisterBindType("mercury:203.2TD", version3)
	boggart.RegisterBindType("mercury:230", version3)
	boggart.RegisterBindType("mercury:231", version3)
	boggart.RegisterBindType("mercury:233", version3)
	boggart.RegisterBindType("mercury:234", version3)
}
