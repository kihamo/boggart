package openhab

import (
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/openhab/client"
)

type Client struct {
	*client.OpenHABREST
}

func New(address *url.URL, debug bool, logger logger.Logger) *Client {
	cfg := client.DefaultTransportConfig().
		WithSchemes([]string{address.Scheme}).
		WithHost(address.Host)

	cl := &Client{
		OpenHABREST: client.NewHTTPClientWithConfig(nil, cfg),
	}

	if rt, ok := cl.OpenHABREST.Transport.(*httptransport.Runtime); ok {
		username := address.User.Username()
		password, _ := address.User.Password()

		if username != "" && password != "" {
			rt.DefaultAuthentication = httptransport.BasicAuth(username, password)
		}

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}
