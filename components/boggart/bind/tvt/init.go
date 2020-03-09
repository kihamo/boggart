package tvt

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("tvt", Type{})
	boggart.RegisterBindType("praxis:vdr", Type{})
}
