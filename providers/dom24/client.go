package dom24

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/dom24/client"
	"github.com/kihamo/boggart/providers/dom24/client/auth"
	"github.com/kihamo/boggart/providers/dom24/models"
)

const (
	headerClient         = "client"
	headerClientValue    = "Android"
	headerUserAgentValue = "RestSharp/106.11.2.0"
	headerToken          = "acx"
)

type Client struct {
	*client.Dom24

	tokenLock sync.RWMutex
	token     string
}

func New(debug bool, logger logger.Logger) *Client {
	cl := &Client{
		Dom24: client.NewHTTPClient(nil),
	}

	if rt, ok := cl.Dom24.Transport.(*httptransport.Runtime); ok {
		originTransport := rt.Transport

		rt.Transport = swagger.RoundTripperFunc(func(request *http.Request) (response *http.Response, err error) {
			request.Header.Set("User-Agent", headerUserAgentValue)
			request.Header[headerClient] = []string{headerClientValue}

			response, err = originTransport.RoundTrip(request)
			if err != nil {
				return response, err
			}

			// эмулируем http ошибку, так как у них всегда 200 OK
			body, err := ioutil.ReadAll(response.Body)
			if err == nil {
				var (
					errorResponse models.Error
					newBody       []byte
				)

				if err := json.Unmarshal(body, &errorResponse); err == nil {
					if errorResponse.Error != "" {
						if b, err := json.Marshal(errorResponse); err == nil {
							newBody = b
						}
					}
				}

				if len(newBody) > 0 {
					response.StatusCode = http.StatusBadRequest
					response.Status = http.StatusText(http.StatusBadRequest)
					response.Body = ioutil.NopCloser(bytes.NewReader(newBody))
					response.ContentLength = int64(len(newBody))
				} else {
					response.Body = ioutil.NopCloser(bytes.NewReader(body))
				}
			}

			for _, cookie := range response.Cookies() {
				if cookie.Name == headerToken {
					cl.tokenLock.Lock()
					cl.token, _ = url.QueryUnescape(cookie.Value)
					cl.tokenLock.Unlock()
				}
			}

			return response, err
		})

		rt.DefaultAuthentication = runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) (err error) {
			cl.tokenLock.RLock()
			defer cl.tokenLock.RUnlock()

			if cl.token != "" {
				// какие-то конченные разработчики, которые дрочат на регистр в заголовках
				req.GetHeaderParams()[headerToken] = []string{cl.token}
			}

			return nil
		})

		rt.Consumers["text/html"] = runtime.XMLConsumer()
		rt.Consumers["text/xml"] = runtime.XMLConsumer()

		// для bills
		rt.Consumers["image/jpeg"] = runtime.ByteStreamConsumer()
		rt.Consumers["application/pdf"] = runtime.ByteStreamConsumer()

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}

func (c *Client) Login(ctx context.Context, phone, password string) (*models.Account, error) {
	params := auth.NewLoginParamsWithContext(ctx)
	params.Request.Phone = phone
	params.Request.Password = password

	resp, err := c.Dom24.Auth.Login(params)
	if err != nil {
		return nil, err
	}

	return resp.GetPayload(), err
}
