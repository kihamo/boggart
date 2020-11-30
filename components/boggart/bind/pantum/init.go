package pantum

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	t := Type{}
	boggart.RegisterBindType("pantum:m6700", t)
	boggart.RegisterBindType("pantum:m6800", t)
	boggart.RegisterBindType("pantum:m7100", t)
	boggart.RegisterBindType("pantum:m7200", t)
}
