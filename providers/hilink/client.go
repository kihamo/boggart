package hilink

// http://forum.jdtech.pl/Watek-hilink-api-dla-urzadzen-huawei
// https://github.com/knq/hilink

import (
	"context"
	"net"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/providers/hilink/client"
	"github.com/kihamo/boggart/providers/hilink/client/web_server"
)

const (
	headerToken   = "__RequestVerificationToken"
	headerSession = "Cookie"
)

type Client struct {
	*client.HiLink

	runtime *Runtime
}

func New(address string, debug bool, logger logger.Logger) *Client {
	if host, port, err := net.SplitHostPort(address); err == nil && port == "80" {
		address = host
	}

	cl := &Client{
		HiLink: client.NewHTTPClientWithConfig(nil, client.DefaultTransportConfig().WithHost(address)),
	}

	if rt, ok := cl.HiLink.Transport.(*httptransport.Runtime); ok {
		cl.runtime = newRuntime(rt)

		// автоматическая авторизация при первом запросе
		cl.runtime.SetDefaultAuthentication(runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, rg strfmt.Registry) error {
			if r.GetPath() == "/webserver/SesTokInfo" {
				return nil
			}

			if err := cl.Auth(context.Background()); err != nil {
				return err
			}

			return cl.HiLink.Transport.(*httptransport.Runtime).
				DefaultAuthentication.
				AuthenticateRequest(r, rg)
		}))

		rt.Transport = RoundTripper{
			original: rt.Transport,
			runtime:  cl.runtime,
		}

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

	return cl
}

func (c *Client) Auth(ctx context.Context) error {
	params := web_server.NewGetWebServerSessionParams().
		WithContext(ctx)

	response, err := c.WebServer.GetWebServerSession(params)
	if err != nil {
		return err
	}

	token := response.Payload.Token
	session := response.Payload.Session

	c.runtime.SetDefaultAuthentication(runtime.ClientAuthInfoWriterFunc(func(request runtime.ClientRequest, _ strfmt.Registry) error {
		if err := request.SetHeaderParam(headerToken, token); err != nil {
			return err
		}

		return request.SetHeaderParam(headerSession, session)
	}))

	return nil
}
