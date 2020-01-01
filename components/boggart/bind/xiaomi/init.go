package xiaomi

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/xiaomi/roborock/miio"
	"github.com/kihamo/boggart/components/boggart/bind/xiaomi/roborock/root"
	"github.com/kihamo/boggart/components/boggart/bind/xiaomi/scale"
)

func init() {
	boggart.RegisterBindType("xiaomi:roborock:root", root.Type{})
	boggart.RegisterBindType("xiaomi:roborock:miio", miio.Type{})
	boggart.RegisterBindType("xiaomi:scale", scale.Type{})
}
