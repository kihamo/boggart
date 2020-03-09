package tvt

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/kihamo/boggart/providers/tvt/models"
)

type RoundTripper struct {
	original http.RoundTripper
}

func (rt RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	inBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	newInBody := []byte(`<?xml version="1.0" encoding="utf-8" ?><request version="1.0" systemType="NVMS-9000" clientType="WEB">`)
	newInBody = append(newInBody, inBody...)

	req.Body = ioutil.NopCloser(bytes.NewReader(newInBody))
	req.ContentLength = int64(len(newInBody))

	response, err := rt.original.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// эмулируем http ошибку, так как всегда 200 OK
	outBody, err := ioutil.ReadAll(response.Body)
	if err == nil {
		var errorResponse models.Error

		if err := xml.Unmarshal(outBody, &errorResponse); err == nil {
			if errorResponse.ErrorCode > 0 {
				response.StatusCode = http.StatusBadRequest
				response.Status = http.StatusText(http.StatusBadRequest)
			}
		}

		response.Body = ioutil.NopCloser(bytes.NewReader(outBody))
	}

	return response, err
}
