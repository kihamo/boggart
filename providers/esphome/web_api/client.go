package web_api

// https://esphome.io/web-api/index.html

import (
	"net"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/esphome/web_api/client"
)

type Client struct {
	*client.ESPHome
}

func New(address string, debug bool, logger logger.Logger) *Client {
	if host, port, err := net.SplitHostPort(address); err == nil && port == "80" {
		address = host
	}

	cl := &Client{
		ESPHome: client.NewHTTPClientWithConfig(nil, client.DefaultTransportConfig().WithHost(address)),
	}

	if rt, ok := cl.ESPHome.Transport.(*httptransport.Runtime); ok {
		rt.Consumers["text/json"] = runtime.JSONConsumer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}
