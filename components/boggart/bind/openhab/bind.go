package openhab

import (
	"net/http/httputil"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/openhab"
)

type Bind struct {
	di.MQTTBind
	di.MetaBind
	di.LoggerBind
	di.ProbesBind
	di.WidgetBind

	config   *Config
	provider *openhab.Client
	proxy    *httputil.ReverseProxy
}
