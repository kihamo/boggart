package pantum

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("pantum", Type{}, "pantum:m6700", "pantum:m6800", "pantum:m7100", "pantum:m7200")
}
