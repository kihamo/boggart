package hilink

// http://forum.jdtech.pl/Watek-hilink-api-dla-urzadzen-huawei
// https://github.com/knq/hilink

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/kihamo/boggart/providers/hilink/client"
	"github.com/kihamo/boggart/providers/hilink/client/user"
	"github.com/kihamo/boggart/providers/hilink/client/web_server"
)

const (
	headerToken   = "__RequestVerificationToken"
	headerSession = "Cookie"

	TimeFormat = "2006-01-02 15:04:05"
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
		cl.runtime = newRuntime(rt, cl.Auth)

		rt.Transport = RoundTripper{
			original: rt.Transport,
			runtime:  cl.runtime,
		}

		rt.Consumers["text/html"] = runtime.XMLConsumer()
		rt.Consumers["text/xml"] = runtime.XMLConsumer()

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
	params := web_server.NewGetWebServerSessionParamsWithContext(ctx)
	response, err := c.WebServer.GetWebServerSession(params)

	if err != nil {
		return err
	}

	c.runtime.SetAuthenticationLogged(response.Payload.Token, response.Payload.Session)

	return nil
}

func (c *Client) Login(ctx context.Context, username, password string) error {
	paramsAuth := web_server.NewGetWebServerSessionParamsWithContext(ctx)
	responseAuth, err := c.WebServer.GetWebServerSession(paramsAuth)

	if err != nil {
		return err
	}

	c.runtime.SetAuthenticationLogged(responseAuth.Payload.Token, responseAuth.Payload.Session)

	h := sha256.New()
	h.Write([]byte(password))

	passwordHash := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(h.Sum(nil))))

	h.Reset()
	h.Write([]byte(username))
	h.Write([]byte(passwordHash))
	h.Write([]byte(responseAuth.Payload.Token))

	paramsLogin := user.NewLoginParamsWithContext(ctx)
	paramsLogin.Request.Username = username
	paramsLogin.Request.Password = base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(h.Sum(nil))))
	paramsLogin.Request.PasswordType = 4

	_, err = c.User.Login(paramsLogin)

	return err
}
