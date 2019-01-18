package xiaomi

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("xiaomi:roborock:root", RoborockRootType{})
}
