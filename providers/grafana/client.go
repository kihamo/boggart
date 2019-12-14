package grafana

import (
	"net"
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/grafana/client"
)

type Client struct {
	*client.Grafana
}

func New(address string, debug bool, logger logger.Logger) *Client {
	u, _ := url.Parse(address)
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	if u.Port() == "" {
		switch u.Scheme {
		case "https":
			u.Host = net.JoinHostPort(u.Hostname(), "443")
		default:
			u.Host = net.JoinHostPort(u.Hostname(), "80")
		}
	}

	cfg := client.DefaultTransportConfig().WithHost(u.Host).
		WithSchemes([]string{u.Scheme})
	cl := client.NewHTTPClientWithConfig(nil, cfg)

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		if u.User != nil {
			if password, ok := u.User.Password(); ok {
				if u.User.Username() == "api_key" {
					rt.DefaultAuthentication = httptransport.BearerToken(password)
				} else {
					rt.DefaultAuthentication = httptransport.BasicAuth(u.User.Username(), password)
				}
			}
		}

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return &Client{
		Grafana: cl,
	}
}
