package hilink

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/web_server"
)

const (
	headerToken   = "__RequestVerificationToken"
	headerSession = "Cookie"
)

type Client struct {
	*client.HiLink
}

type XMLRoundTripper struct {
	original http.RoundTripper
	runtime  *httptransport.Runtime

	mutex sync.Mutex
}

func New(address string, debug bool, logger logger.Logger) *Client {
	if host, port, err := net.SplitHostPort(address); err == nil && port == "80" {
		address = host
	}

	cfg := client.DefaultTransportConfig().WithHost(address)
	cl := client.NewHTTPClientWithConfig(nil, cfg)

	if rt, ok := cl.Transport.(*httptransport.Runtime); ok {
		rt.Transport = XMLRoundTripper{
			original: rt.Transport,
			runtime:  rt,
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
			if err = r.SetHeaderParam(headerToken, response.Payload.Token); err != nil {
				return err
			}

			return r.SetHeaderParam(headerSession, response.Payload.Session)
		})
	}

	return nil
}

func (rt XMLRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	/*
		if req.Method == http.MethodPost && req.ContentLength > 0 && req.Header.Get("Content-Type") == "application/xml" {
			//req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}

			bytes.Index(body, []byte(">"))

			start := bytes.Index(body, []byte(">"))
			stop := bytes.LastIndex(body, []byte("/"))

			if start > -1 && stop > -1 {
				body = append([]byte(`<?xml version: "1.0" encoding="UTF-8"?><request`), body[start:stop+1]...)
				body = append(body, []byte("request>")...)
			}

			req.Body = ioutil.NopCloser(bytes.NewReader(body))
			req.ContentLength = int64(len(body))
		}
	*/

	response, err := rt.original.RoundTrip(req)
	if err == nil {
		token := response.Header.Get(headerToken)
		if token != "" {
			rt.mutex.Lock()

			rt.runtime.DefaultAuthentication = runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
				if err = r.SetHeaderParam(headerToken, token); err != nil {
					return err
				}

				return r.SetHeaderParam(headerSession, req.Header.Get(headerSession))
			})

			rt.mutex.Unlock()
		}
	}

	return response, err
}
