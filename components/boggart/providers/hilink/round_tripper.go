package hilink

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type RoundTripper struct {
	original http.RoundTripper
	runtime  *Runtime
}

func (rt RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
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
			rt.runtime.SetDefaultAuthentication(runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) (err error) {
				err = r.SetHeaderParam(headerToken, token)
				if err == nil {
					err = r.SetHeaderParam(headerSession, req.Header.Get(headerSession))
				}

				return err
			}))
		}
	}

	return response, err
}
