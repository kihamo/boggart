package wiim

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/wiim/httpapi/client"
)

// Support Linkplay
// http:// :49152/description.xml
// https://github.com/n4archive/LinkPlayAPI/blob/master/api.md

// https://developer.arylic.com/httpapi/#http-api

type Client struct {
	*client.Wiim
}

func New(address *url.URL, debug bool, logger logger.Logger) *Client {
	cfg := client.DefaultTransportConfig().
		WithSchemes([]string{address.Scheme}).
		WithHost(address.Host)

	cl := &Client{
		Wiim: client.NewHTTPClientWithConfig(nil, cfg),
	}

	rt, ok := cl.Transport.(*httptransport.Runtime)
	if !ok {
		return cl
	}

	rt.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	rt.Consumers["text/html"] = runtime.JSONConsumer()
	rt.Consumers["text/plain"] = runtime.TextConsumer()

	if logger != nil {
		rt.SetLogger(logger)
	}

	rt.SetDebug(debug)

	cl.SetTransport(&transport{
		proxied: rt,
	})

	rt.Transport = &roundTripper{
		proxied: rt.Transport,
	}

	return cl
}
