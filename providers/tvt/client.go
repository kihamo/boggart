package tvt

import (
	"net"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/tvt/client"
)

type Client struct {
	*client.TVT
}

func New(address string, user, password string, debug bool, logger logger.Logger) *Client {
	if host, port, err := net.SplitHostPort(address); err == nil && port == "80" {
		address = host
	}

	cfg := client.DefaultTransportConfig().WithHost(address)
	cl := client.NewHTTPClientWithConfig(nil, cfg)

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		rt.DefaultAuthentication = httptransport.BasicAuth(user, password)
		rt.Transport = RoundTripper{
			original: rt.Transport,
		}

		rt.Consumers["text/xml"] = runtime.XMLConsumer()
		rt.Producers["text/xml"] = runtime.XMLProducer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return &Client{
		TVT: cl,
	}
}
