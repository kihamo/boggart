package device

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"net"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/myheat/device/client"
)

const (
	DefaultUsername = "myheat"
	DefaultPassword = "myheat"
)

type Client struct {
	*client.MyHeat
}

func New(address string, debug bool, logger logger.Logger) *Client {
	if host, port, err := net.SplitHostPort(address); err == nil && port == "80" {
		address = host
	}

	cfg := client.DefaultTransportConfig().WithHost(address)
	cl := client.NewHTTPClientWithConfig(nil, cfg)

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		rt.DefaultAuthentication = httptransport.BasicAuth(DefaultUsername, DefaultPassword)

		fmt.Println(rt.Consumers)
		fmt.Println(rt.Producers)

		rt.Consumers["text/json"] = runtime.JSONConsumer()
		rt.Producers["text/json"] = runtime.JSONProducer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return &Client{
		MyHeat: cl,
	}
}
