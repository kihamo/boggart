package openhab

import (
	"net/http"
	"net/http/httputil"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/openhab"
)

type Type struct{}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	config := c.(*Config)

	bind := &Bind{
		config: config,
	}

	l := swagger.NewLogger(
		func(message string) {
			bind.Logger().Info(message)
		},
		func(message string) {
			bind.Logger().Debug(message)
		})

	bind.provider = openhab.New(&config.Address.URL, config.Debug, l)

	bind.proxy = httputil.NewSingleHostReverseProxy(&config.Address.URL)
	director := bind.proxy.Director

	bind.proxy.Director = func(request *http.Request) {
		request.Host = config.Address.Host

		if username := config.Address.User.Username(); username != "" {
			password, _ := config.Address.User.Password()
			request.SetBasicAuth(username, password)
		}

		director(request)
	}

	if config.ProxyEnabled {
		bind.proxyServer = &http.Server{
			Addr: config.ProxyAddress,
		}
	}

	return bind, nil
}
