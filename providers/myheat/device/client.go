package device

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"sync/atomic"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/myheat/device/client"
	"github.com/kihamo/boggart/providers/myheat/device/client/auth"
	"github.com/kihamo/boggart/providers/myheat/device/models"
)

const (
	DefaultUsername = "myheat"
	DefaultPassword = "myheat"

	cookieId = "EspSessId"
)

type Client struct {
	*client.MyHeat

	authLock atomic.Bool
}

func New(address string, debug bool, logger logger.Logger) *Client {
	if host, port, err := net.SplitHostPort(address); err == nil && port == "80" {
		address = host
	}

	cfg := client.DefaultTransportConfig().WithHost(address)
	transport := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	transport.Jar, _ = cookiejar.New(nil)

	cl := &Client{
		MyHeat: client.New(transport, nil),
	}

	originTransport := transport.Transport
	transport.Transport = swagger.RoundTripperFunc(func(request *http.Request) (response *http.Response, err error) {
		response, err = originTransport.RoundTrip(request)
		if err == nil {
			for _, v := range response.Cookies() {
				if v.Name == cookieId {
					cl.authLock.Store(true)
					break
				}
			}
		}

		return response, err
	})

	transport.DefaultAuthentication = runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) (err error) {
		if !cl.authLock.Load() && req.GetPath() != "/login" {
			params := auth.NewLoginParams().
				WithRequest(&models.LoginRequest{
					Login:    DefaultUsername,
					Password: DefaultPassword,
				})

			_, err = cl.Auth.Login(params)

			return err
		}

		return nil
	})

	transport.Consumers["text/json"] = runtime.JSONConsumer()
	transport.Producers["text/json"] = runtime.JSONProducer()

	if logger != nil {
		transport.SetLogger(logger)
	}

	transport.SetDebug(debug)

	return cl
}
