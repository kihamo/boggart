package wiim

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/wiim/client"
)

// Support Linkplay
// http:// :49152/description.xml
// https://github.com/n4archive/LinkPlayAPI/blob/master/api.md

// https://developer.arylic.com/httpapi/#http-api
// getStaticIP -- not support
// getsyslog -- support
// getPlayerStatus -- support
// GetTrackNumber -- not support
// getLocalPlayList -- not support
// multiroom:getSlaveList -- support

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

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		rt.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		rt.Consumers["text/html"] = runtime.JSONConsumer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}
