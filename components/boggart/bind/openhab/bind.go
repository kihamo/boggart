package openhab

import (
	"net/http"
	"net/http/httputil"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/openhab"
	m "github.com/kihamo/shadow/components/logging/http"
)

type Bind struct {
	di.MQTTBind
	di.MetaBind
	di.LoggerBind
	di.ProbesBind
	di.WidgetBind

	config      *Config
	provider    *openhab.Client
	proxy       *httputil.ReverseProxy
	proxyServer *http.Server
}

func (b *Bind) Run() error {
	if b.config.ProxyEnabled {
		mw := m.ServerMiddleware(b.Logger())
		b.proxyServer.Handler = mw(http.HandlerFunc(b.proxyHandler))

		go func() {
			if err := b.proxyServer.ListenAndServe(); err != nil {
				b.Logger().Error("Failed serve with error " + err.Error())
			}
		}()
	}

	return nil
}

func (b *Bind) Close() error {
	if b.config.ProxyEnabled {
		return b.proxyServer.Close()
	}

	return nil
}

func (b *Bind) proxyHandler(rw http.ResponseWriter, req *http.Request) {
	b.proxy.ServeHTTP(rw, req)
}
