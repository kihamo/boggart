package openweathermap

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterBindType("openweathermap", Type{})
}
