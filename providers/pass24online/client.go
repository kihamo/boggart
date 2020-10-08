package pass24online

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/logger"
	"github.com/go-openapi/strfmt"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/providers/pass24online/client"
	"github.com/kihamo/boggart/providers/pass24online/client/auth"
	"github.com/kihamo/boggart/providers/pass24online/models"
)

type Client struct {
	*client.Pass24online

	transport http.RoundTripper

	phone    string
	password string

	token     string
	tokenOnce *atomic.Once
}

func New(phone, password string, debug bool, logger logger.Logger) *Client {
	cl := &Client{
		Pass24online: client.NewHTTPClient(nil),
		phone:        phone,
		password:     password,
		tokenOnce:    &atomic.Once{},
	}

	if rt, ok := cl.Pass24online.Transport.(*httptransport.Runtime); ok {
		cl.transport = rt.Transport

		rt.Transport = cl
		rt.DefaultAuthentication = cl

		if logger != nil {
			rt.SetLogger(logger)
		}

		rt.SetDebug(debug)
	}

	return cl
}

func (c *Client) RoundTrip(req *http.Request) (*http.Response, error) {
	response, err := c.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// эмулируем http ошибку, так как у апи всегда 200 OK
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		response.Body = ioutil.NopCloser(bytes.NewReader(body))
		response.ContentLength = int64(len(body))

		var errorResponse models.ErrorResponse

		if err := json.Unmarshal(body, &errorResponse); err == nil && errorResponse.Error != nil {
			switch errorResponse.Error.Code {
			case "UNAUTHORIZED":
				response.StatusCode = http.StatusUnauthorized
			case "UNAUTHENTICATED":
				response.StatusCode = http.StatusForbidden

				// возможно авторизация слетела, поэтому скидываем внутренний токен
				c.tokenOnce.Reset()
			default:
				response.StatusCode = http.StatusBadRequest
			}

			response.Status = http.StatusText(response.StatusCode)
		}
	}

	return response, err
}

func (c *Client) AuthenticateRequest(req runtime.ClientRequest, reg strfmt.Registry) (err error) {
	if req.GetPath() == "/auth/login/" {
		return nil
	}

	c.tokenOnce.Do(func() {
		var response *auth.LoginOK

		params := auth.NewLoginParams().
			WithPhone(c.phone).
			WithPassword(c.password)

		response, err = c.Pass24online.Auth.Login(params)

		if err == nil {
			c.token = response.GetPayload().Body
		}
	})

	if err != nil {
		c.tokenOnce.Reset()

		return err
	}

	req.SetHeaderParam("Authorization", "Bearer "+c.token)

	return nil
}
