package smcenter

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/smcenter/client"
	"github.com/kihamo/boggart/providers/smcenter/client/auth"
	"github.com/kihamo/boggart/providers/smcenter/models"
)

const (
	headerClient         = "client"
	headerClientValue    = "Android"
	headerUserAgentValue = "RestSharp/106.11.2.0"
	headerToken          = "acx"
)

type Client struct {
	*client.SMCenter

	tokenLock    sync.RWMutex
	tokenValue   string
	tokenExpires *time.Time
}

func New(basePath, phone, password string, debug bool, logger logger.Logger) *Client {
	cfg := client.DefaultTransportConfig()

	if basePath != "" {
		cfg.BasePath = basePath
	}

	transport := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)

	cl := &Client{
		SMCenter: client.New(transport, nil),
	}

	originTransport := transport.Transport

	transport.Transport = swagger.RoundTripperFunc(func(request *http.Request) (response *http.Response, err error) {
		request.Header.Set("User-Agent", headerUserAgentValue)
		request.Header[headerClient] = []string{headerClientValue}

		response, err = originTransport.RoundTrip(request)
		if err != nil {
			return response, err
		}

		if response.StatusCode == http.StatusUnauthorized {
			cl.setToken("", nil)

			return response, err
		}

		if response.StatusCode != http.StatusOK {
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
				v, _ := url.QueryUnescape(cookie.Value)
				cl.setToken(v, &cookie.Expires)
			}
		}

		return response, err
	})

	transport.DefaultAuthentication = runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) (err error) {
		token := cl.token()

		if token == "" && req.GetPath() != "/auth/login" {
			params := auth.NewLoginParams()
			params.Request.Phone = phone
			params.Request.Password = password

			if _, err := cl.SMCenter.Auth.Login(params); err == nil {
				token = cl.token()
			}
		}

		if token != "" {
			// какие-то конченные разработчики, которые дрочат на регистр в заголовках
			req.GetHeaderParams()[headerToken] = []string{token}
		}

		return nil
	})

	transport.Consumers["text/html"] = runtime.XMLConsumer()
	transport.Consumers["text/xml"] = runtime.XMLConsumer()

	// для bills
	transport.Consumers["image/jpeg"] = runtime.ByteStreamConsumer()
	transport.Consumers["image/png"] = runtime.ByteStreamConsumer()
	transport.Consumers["application/pdf"] = runtime.ByteStreamConsumer()

	// в SetLogger и SetDebug race, так как он еще и ставится глобально, what???
	if logger != nil {
		swagger.SetLogger(transport, logger)
	}

	swagger.SetDebug(transport, debug)

	return cl
}

func (c *Client) setToken(token string, expires *time.Time) {
	if token == "" {
		expires = nil
	}

	c.tokenLock.Lock()
	c.tokenValue = token
	c.tokenExpires = expires
	c.tokenLock.Unlock()
}

func (c *Client) token() string {
	c.tokenLock.RLock()
	value := c.tokenValue
	expires := c.tokenExpires
	c.tokenLock.RUnlock()

	if expires != nil && !expires.IsZero() && time.Now().After(*expires) {
		value = ""

		c.tokenLock.Lock()
		c.tokenValue = ""
		c.tokenExpires = nil
		c.tokenLock.Unlock()
	}

	return value
}

func (c *Client) Login(ctx context.Context, phone, password string) (*models.Account, error) {
	params := auth.NewLoginParamsWithContext(ctx)
	params.Request.Phone = phone
	params.Request.Password = password

	resp, err := c.SMCenter.Auth.Login(params)
	if err != nil {
		return nil, err
	}

	return resp.GetPayload(), err
}
