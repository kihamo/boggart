package hilink

import (
	"context"
	"net"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/web_server"
)

const (
	TokenHeader = "__RequestVerificationToken"
)

type Client struct {
	*client.HiLink
}

func New(address string, debug bool, logger logger.Logger) *Client {
	if host, port, err := net.SplitHostPort(address); err == nil && port == "80" {
		address = host
	}

	cfg := client.DefaultTransportConfig().WithHost(address)
	cl := client.NewHTTPClientWithConfig(nil, cfg)

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		rt.Consumers["text/html"] = runtime.XMLConsumer()

		// что бы скачивались файлы с изображениями
		rt.Consumers["image/jpeg"] = runtime.ByteStreamConsumer()
		rt.Consumers["image/png"] = runtime.ByteStreamConsumer()
		rt.Consumers["image/gif"] = runtime.ByteStreamConsumer()

		// для alert stream
		rt.Consumers["multipart/mixed"] = runtime.ByteStreamConsumer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return &Client{
		HiLink: cl,
	}
}

func (c *Client) Auth(ctx context.Context) error {
	params := web_server.NewGetWebServerSessionParams().
		WithContext(ctx)

	response, err := c.WebServer.GetWebServerSession(params)
	if err != nil {
		return err
	}

	if rt, ok := c.Transport.(*httptransport.Runtime); ok {
		rt.DefaultAuthentication = runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
			if err = r.SetHeaderParam(TokenHeader, response.Payload.Token); err != nil {
				return err
			}

			return r.SetHeaderParam("Cookie", response.Payload.Session)
		})
	}

	return nil
}
