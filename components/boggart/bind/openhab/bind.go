package openhab

import (
	"context"
	"net/http"
	"net/http/httputil"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/openhab3"
	m "github.com/kihamo/shadow/components/logging/http"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind

	provider    *openhab3.Client
	proxy       *httputil.ReverseProxy
	proxyServer *http.Server
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() (err error) {
	cfg := b.config()
	opts := make([]openhab3.ClientOption, 0, 2)

	if username := cfg.Address.User.Username(); username != "" {
		password, _ := cfg.Address.User.Password()
		var autProvider *securityprovider.SecurityProviderApiKey

		autProvider, err = securityprovider.NewSecurityProviderApiKey("header", username, password)
		if err != nil {
			return err
		}

		opts = append(opts, openhab3.WithRequestEditorFn(autProvider.Intercept))
	}

	if cfg.Debug {
		opts = append(opts, openhab3.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			dump, err := httputil.DumpRequestOut(req, true)

			if err != nil {
				return err
			}

			b.Logger().Debugf("\n\n%q", dump)

			return nil
		}))
	}

	b.provider, err = openhab3.NewClient(cfg.Address.URL.String()+"/rest", opts...)
	if err != nil {
		return err
	}

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

	b.Meta().SetLink(&cfg.Address.URL)

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
