package openhab

import (
	"net/http"
	"net/http/httputil"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/openhab"
	m "github.com/kihamo/shadow/components/logging/http"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind

	provider    *openhab.Client
	proxy       *httputil.ReverseProxy
	proxyServer *http.Server
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	b.provider = openhab.New(&cfg.Address.URL, cfg.Debug, swagger.NewLogger(
		func(message string) {
			b.Logger().Info(message)
		},
		func(message string) {
			b.Logger().Debug(message)
		}))

	b.proxy = httputil.NewSingleHostReverseProxy(&cfg.Address.URL)
	director := b.proxy.Director

	b.proxy.Director = func(request *http.Request) {
		request.Host = cfg.Address.Host

		if username := cfg.Address.User.Username(); username != "" {
			password, _ := cfg.Address.User.Password()
			request.SetBasicAuth(username, password)
		}

		director(request)
	}

	if cfg.ProxyEnabled {
		b.proxyServer = &http.Server{
			Addr: cfg.ProxyAddress,
		}

		mw := m.ServerMiddleware(b.Logger())
		b.proxyServer.Handler = mw(http.HandlerFunc(b.proxyHandler))

		go func() {
			if err := b.proxyServer.ListenAndServe(); err != nil {
				b.Logger().Error("Failed serve with error " + err.Error())
			}
		}()
	} else {
		b.proxyServer = nil
	}

	return nil
}

func (b *Bind) Close() error {
	if b.config().ProxyEnabled {
		return b.proxyServer.Close()
	}

	return nil
}

func (b *Bind) proxyHandler(rw http.ResponseWriter, req *http.Request) {
	b.proxy.ServeHTTP(rw, req)
}
