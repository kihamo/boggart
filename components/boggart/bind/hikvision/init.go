package hikvision

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/types"
)

var defaultNTPServer types.URL

func init() {
	defaultNTPServer, _ = types.ParseURL("ntp://time.windows.com:123")

	boggart.RegisterBindType("hikvision", Type{})
}
