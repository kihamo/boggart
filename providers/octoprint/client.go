package octoprint

import (
	"net"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/octoprint/client"
)

type Client struct {
	*client.OctoPrint
}

func New(address, apiKey string, debug bool, logger logger.Logger) *Client {
	if host, port, err := net.SplitHostPort(address); err == nil && port == "80" {
		address = host
	}

	cl := &Client{
		OctoPrint: client.NewHTTPClientWithConfig(nil, client.DefaultTransportConfig().WithHost(address)),
	}

	if rt, ok := cl.OctoPrint.Transport.(*httptransport.Runtime); ok {
		rt.DefaultAuthentication = httptransport.APIKeyAuth("X-Api-Key", "header", apiKey)
		rt.Consumers["image/png"] = runtime.ByteStreamConsumer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}
