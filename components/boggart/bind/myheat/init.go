package myheat

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("myheat", Type{}, "myheat:smart2")
}
