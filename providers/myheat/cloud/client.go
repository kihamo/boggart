package cloud

import (
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/myheat/cloud/client"
)

type Client struct {
	*client.MyHeatCloud
}

func New(login, key string, debug bool, logger logger.Logger) *Client {
	cl := client.NewHTTPClient(nil)

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)

		cl.SetTransport(&transport{
			proxied: rt,
			login:   login,
			key:     key,
		})

		rt.Transport = &roundTripper{
			proxied: rt.Transport,
		}
	}

	return &Client{
		MyHeatCloud: cl,
	}
}
