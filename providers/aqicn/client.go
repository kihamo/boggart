package aqicn

import (
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/aqicn/client"
)

type Client struct {
	*client.AirQuality
}

func New(token string, debug bool, logger logger.Logger) *Client {
	cl := &Client{
		AirQuality: client.Default,
	}

	if rt, ok := cl.AirQuality.Transport.(*httptransport.Runtime); ok {
		rt.DefaultAuthentication = httptransport.APIKeyAuth("token", "query", token)

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}

func Icon(id string) *url.URL {
	u := &url.URL{}
	u.Host = "aqicn.org"
	u.Path = "/images/feeds/" + id

	return u
}
