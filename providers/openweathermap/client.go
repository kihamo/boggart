package openweathermap

import (
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/openweathermap/client"
)

type Client struct {
	*client.OpenWeather
}

func New(apiKey string, debug bool, logger logger.Logger) *Client {
	cl := &Client{
		OpenWeather: client.NewHTTPClient(nil),
	}

	if rt, ok := cl.OpenWeather.Transport.(*httptransport.Runtime); ok {
		rt.DefaultAuthentication = httptransport.APIKeyAuth("appid", "query", apiKey)

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}
