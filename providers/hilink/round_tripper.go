package hilink

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kihamo/boggart/providers/hilink/models"
)

type RoundTripper struct {
	original http.RoundTripper
	runtime  *Runtime
}

func (rt RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	response, err := rt.original.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if response.ContentLength > 0 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var apiError models.Error

		if err := xml.Unmarshal(body, &apiError); err == nil {
			switch apiError.Code {
			case ErrorTokenWrong, ErrorSessionWrong, ErrorSessionTokenWrong:
				rt.runtime.SetAuthenticationAnonymous()
			}
		}

		response.Body = ioutil.NopCloser(bytes.NewReader(body))
		response.ContentLength = int64(len(body))
	}

	token := response.Header.Get(headerToken)
	session := req.Header.Get(headerSession)

	if token != "" && session != "" {
		rt.runtime.SetAuthenticationLogged(token, session)
	}

	// эмулируем http ошибку, так как у hilink всегда 200 OK
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		var (
			errorResponse models.Error
			newBody       []byte
		)

		if err := xml.Unmarshal(body, &errorResponse); err == nil {
			if errorResponse.Code > 0 {
				if len(errorResponse.Message) == 0 {
					errorResponse.Message = ErrorMessage(errorResponse.Code)
				}

				if b, err := xml.Marshal(errorResponse); err == nil {
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

	// go-swagger на каждый реквест создает нового клиента, поэтому организуем
	// прокидывание печенек самостоятельно
	cookies := response.Cookies()

	if len(cookies) > 0 {
		u, _ := url.Parse(req.URL.Scheme + "://" + req.URL.Host + rt.runtime.BasePath)
		rt.runtime.Jar.SetCookies(u, response.Cookies())
	}

	return response, err
}
